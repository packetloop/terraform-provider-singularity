package main

import (
	"github.com/hashicorp/terraform/plugin"
	singularity "github.com/packetloop/terraform-provider-singularity/mesos_singularity"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: singularity.Provider})
}
