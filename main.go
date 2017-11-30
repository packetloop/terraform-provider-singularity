package main

import (
	"github.com/hashicorp/terraform/plugin"
	singularity "github.com/lenfree/terraform-provider-singularity/mesos_singularity"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: singularity.Provider})
}
