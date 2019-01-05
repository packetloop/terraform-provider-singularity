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
  request_id             = "myrequestphewphew2"
  request_type           = "SCHEDULED"
  schedule               = "0 7 * * *"
  schedule_type          = "CRON"
}

resource "singularity_docker_deploy" "phewphew" {
  deploy_id        = "mydeployphewphewhello4"
  command          = "bash"
  args             = ["-xc", "date"]
  request_id       = "${singularity_request.phewphew.id}"

  docker_info {
    force_pull_image = false
    network          = "BRIDGE"
    image            = "ubuntu"

    port_mapping {
      host_port           = 0
      container_port      = 8888
      container_port_type = "LITERAL"
      host_port_type      = "FROM_OFFER"
      protocol            = "tcp"
    }
  }

  resources {
    cpus      = 2
    memory_mb = 128
  }
}
`
