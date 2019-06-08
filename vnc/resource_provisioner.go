package vnc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform/communicator"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	defaultTimeout  = 5 * time.Minute
	defaultBootWait = time.Duration(0)
)

// Provisioner returns new VNC provisioner
func Provisioner() terraform.ResourceProvisioner {
	return &schema.Provisioner{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "5m",
			},
			"boot_wait": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inline": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				PromoteSingle: true,
				Optional:      true,
				ConflictsWith: []string{"script", "scripts"},
			},
			"script": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"inline", "scripts"},
			},
			"scripts": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"script", "inline"},
			},
		},

		ApplyFunc: applyFn,
	}
}

// Apply executes the remote exec provisioner
func applyFn(ctx context.Context) error {
	data := ctx.Value(schema.ProvConfigDataKey).(*schema.ResourceData)
	o := ctx.Value(schema.ProvOutputKey).(terraform.UIOutput)

	// Get a new vnc communicator
	conn, err := newVNCConnection(data)
	if err != nil {
		return err
	}

	// Collect the scripts
	scripts, err := collectScripts(data)
	if err != nil {
		return err
	}
	for _, s := range scripts {
		defer s.Close()
	}

	waitForBoot(data)

	// Copy and execute each script
	if err := runScripts(ctx, o, conn, scripts); err != nil {
		return err
	}

	return nil
}

func waitForBoot(data *schema.ResourceData) {
	bootWait := safeDuration(data.Get("boot_wait").(string), defaultBootWait)
	if bootWait > 0 {
		log.Println("[DEBUG] Waiting boot...")
		time.Sleep(bootWait)
	}
}

// generateScripts takes the configuration and creates a script from each inline config
func generateScripts(d *schema.ResourceData) ([]string, error) {
	var lines []string
	for _, l := range d.Get("inline").([]interface{}) {
		line, ok := l.(string)
		if !ok {
			return nil, fmt.Errorf("Error parsing %v as a string", l)
		}
		lines = append(lines, line)
	}
	lines = append(lines, "")

	return lines, nil
}

// collectScripts is used to collect all the scripts we need
// to execute in preparation for copying them.
func collectScripts(d *schema.ResourceData) ([]io.ReadCloser, error) {
	// Check if inline
	if _, ok := d.GetOk("inline"); ok {
		scripts, err := generateScripts(d)
		if err != nil {
			return nil, err
		}

		var r []io.ReadCloser
		for _, script := range scripts {
			r = append(r, ioutil.NopCloser(bytes.NewReader([]byte(script))))
		}

		return r, nil
	}

	// Collect scripts
	var scripts []string
	if script, ok := d.GetOk("script"); ok {
		scr, ok := script.(string)
		if !ok {
			return nil, fmt.Errorf("Error parsing script %v as string", script)
		}
		scripts = append(scripts, scr)
	}

	if scriptList, ok := d.GetOk("scripts"); ok {
		for _, script := range scriptList.([]interface{}) {
			scr, ok := script.(string)
			if !ok {
				return nil, fmt.Errorf("Error parsing script %v as string", script)
			}
			scripts = append(scripts, scr)
		}
	}

	// Open all the scripts
	var fhs []io.ReadCloser
	for _, s := range scripts {
		fh, err := os.Open(s)
		if err != nil {
			for _, fh := range fhs {
				fh.Close()
			}
			return nil, fmt.Errorf("Failed to open script '%s': %v", s, err)
		}
		fhs = append(fhs, fh)
	}

	// Done, return the file handles
	return fhs, nil
}

// runScripts is used to copy and execute a set of scripts
func runScripts(
	ctx context.Context,
	o terraform.UIOutput,
	conn *vncConnection,
	scripts []io.ReadCloser) error {

	retryCtx, cancel := context.WithTimeout(ctx, conn.timeout)
	defer cancel()

	// Wait and retry until we establish the connection
	err := communicator.Retry(retryCtx, func() error {
		return conn.Connect(o)
	})
	if err != nil {
		return err
	}

	// Wait for the context to end and then disconnect
	go func() {
		<-ctx.Done()
		conn.Close() // nolint ignore error
	}()

	for _, script := range scripts {
		if err := conn.SendCommand(ctx, script); err != nil {
			return err
		}
	}

	return nil
}
