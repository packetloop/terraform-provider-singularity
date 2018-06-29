package mesos_singularity

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSingularityDockerDeploy_importBasic(t *testing.T) {
	resourceName := "singularity_docker_deploy.phewphew"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityDockerDeployDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSingularityDeployDockerConfigImport,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCheckSingularityDeployDockerConfigImport = `
resource "singularity_request" "phewphew" {
	request_id             = "myrequestphewphew"
	request_type           = "SCHEDULED"
	schedule               = "0 7 * * *"
	schedule_type          = "CRON"
}
resource "singularity_docker_deploy" "phewphew" {
	deploy_id        = "mydeployphewphew3"
	force_pull_image = false
	network          = "BRIDGE"
	image            = "golang:latest"
	cpu              = 2
	memory           = 128
	command          = "bash"
	args             = ["-xc", "date"]
	request_id       = "${singularity_request.phewphew.id}"
}
`
