package mesos_singularity

import (
	"testing"

  "github.com/hashicorp/terraform/helper/resource"
  singularity "github.com/lenfree/go-singularity"

)

func TestGetDeployMD5(t *testing.T) {
  d := singularity.NewDeploy("my").
    SetCommand("bash").
    SetRequestID("test")

  r := singularity.NewDeployRequest().
    AttachDeploy(d).
    Build()

  expected := "cfd8e5a207488d876b4a041d2824d1f9"
  sum := calculateDeployMD5(r)
  if sum != expected {
      t.Errorf("getDeploy(%+v), expected %v, got %v", r, sum, expected)
  }
}
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
  command          = "bash"
  args             = ["-xc", "date"]
  request_id       = "${singularity_request.phewphew.id}"

  container_info{
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
  }

  resources {
    cpus      = 2
    memory_mb = 128
  }
}
`
