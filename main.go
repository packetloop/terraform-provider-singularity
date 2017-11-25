package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/lenfree/terraform-provider-singularity/singularity"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: singularity.Provider})
}
