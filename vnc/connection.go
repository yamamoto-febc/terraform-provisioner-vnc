package vnc

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/go-vnc"
)

type vncConnection struct {
	timeout  time.Duration
	host     string
	port     int
	password string

	nc net.Conn
	vc *vnc.ClientConn
	d  bootcommand.BCDriver
}

func newVNCConnection(data *schema.ResourceData) (*vncConnection, error) {
	host := data.Get("host").(string)
	port := data.Get("port").(int)
	password := data.Get("password").(string)
	timeout := data.Get("timeout").(string)

	log.Printf("[DEBUG] VNC connect to: %s:%d\n", host, port)

	return &vncConnection{
		host:     host,
		port:     port,
		password: password,
		timeout:  safeDuration(timeout, defaultTimeout),
	}, nil
}

// safeDuration returns either the parsed duration or a default value
func safeDuration(dur string, defaultDur time.Duration) time.Duration {
	d, err := time.ParseDuration(dur)
	if err != nil {
		log.Printf("Invalid duration '%s', using default of %s", dur, defaultDur)
		return defaultDur
	}
	return d
}

// Close implements io.Closer
func (v *vncConnection) Close() error {
	closers := []io.Closer{v.vc, v.nc}

	for _, c := range closers {
		if c == nil {
			continue
		}
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Connect opens new VNC connection
func (v *vncConnection) Connect(o terraform.UIOutput) error {
	log.Println("[DEBUG] Opening VNC connection")
	nc, err := net.Dial("tcp", fmt.Sprintf("%s:%d", v.host, v.port))
	if err != nil {
		return err
	}
	v.nc = nc
	log.Println("[DEBUG] VNC connection opened")

	var auth []vnc.ClientAuth
	if len(v.password) > 0 {
		auth = []vnc.ClientAuth{&vnc.PasswordAuth{Password: v.password}}
	} else {
		auth = []vnc.ClientAuth{new(vnc.ClientAuthNone)}
	}

	vncConn, err := vnc.Client(nc, &vnc.ClientConfig{Auth: auth, Exclusive: true})
	if err != nil {
		return err
	}
	v.vc = vncConn
	v.d = newVNCDriver(v.vc, defaultKeyInterval)
	return nil
}

// SendCommand sends command via VNC connection
func (v *vncConnection) SendCommand(ctx context.Context, script io.ReadCloser) error {
	cmd, err := ioutil.ReadAll(script)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] sending vnc command: %s\n", string(cmd))
	seq, err := bootcommand.GenerateExpressionSequence(string(cmd))
	if err != nil {
		return err
	}
	if err := seq.Do(ctx, v.d); err != nil {
		return nil
	}
	return nil
}
