package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/yamamoto-febc/terraform-provisioner-vnc/vnc"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: vnc.Provisioner,
	})
}
