package mesos_singularity

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOST", nil),
				Description: "The Singularity API endpoint to interface with.",
			},

			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PORT", nil),
				Description: "The Singularity Port to connect to.",
			},

			"retry": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("retry", 3),
				Description: "Number of times to retry when Singularity makes http requests. Defaults to 3.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"singularity_request":       resourceRequest(),
			"singularity_docker_deploy": resourceDockerDeploy(),
		},

		/* DataSources placeholder
		DataSourcesMap: map[string]*schema.Resource{
		},
		*/

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host:  d.Get("host").(string),
		Port:  d.Get("port").(int),
		Retry: d.Get("retry").(int),
	}

	return config.Client()
}
