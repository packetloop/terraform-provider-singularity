package mesos_singularity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	singularity "github.com/lenfree/go-singularity"
)

func TestAccSingularityDeploy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityDeployConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityDeployExists("singularity_deploy.foo"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "request_id", "foo-test-id"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "command", "./start.sh"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "resources.cpus", "1"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "resources.memoryMb", "2048"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "containerInfo.type.", "DOCKER"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "env.NODE_REQUEST_TIMEOUT", "2m"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "env.NODE_REQUEST_NODE_PORT", "3000"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "env.NODE_REQUEST_NODE_HOST", "0.0.0.0"),
					resource.TestCheckResourceAttr(
						"singularity_deploy.foo", "env.SERVICE_TAGS.club_name", "myclub"),
				),
			},
		},
	})
}

// create a resource singularity_docker and attach to
// singuarltiy_deploy. Maybe a docker env, port resource?
const testAccCheckSingularityDeployConfig = `
  resource "singularity_deploy" "foo" {
  request_id = "test-deploy"
  command = "/start.sh"

  resources {
    cpus: 1
    memoryMb: 2048
  }

  containerInfo {
	type: DOCKER
  }

  env {
    NODE_REQUEST_TIMEOUT = "2m"
    NODE_PORT = "3000"
    NODE_HOST = "0.0.0.0"
    SERVICE_TAGS {
	 "club_name" = "myclub"
	 "club_id" = "myid"
  }
}
`

func testAccCheckSingularityDeployExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Conn).sclient
		return SingularityRequestExistsHelper(s, client)
	}
}

func testCheckSingularityDeployDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "singularity_deploy" {
			continue
		}

		requestID := res.Primary.ID

		client := testAccProvider.Meta().(*Conn).sclient
		// To check if deploy exists, we query /api/requests/request_id and
		// check if a activedeploy id exists or not. There is no REST
		// endpoint to send GET method to query deploy by id directly.
		data, err := client.GetRequestByID(requestID)
		if err != nil {
			return nil
		}
		// If request_id does not exists, it gets a response status code 404 Not Found.
		if data.RestyResponse.StatusCode() != 404 {
			return fmt.Errorf("There was an error deleting request id '%s'", requestID)
		}
		// If there is a deploy id, a deploy exists.
		if data.Body.ActiveDeploy.ID != "" {
			return fmt.Errorf("There was an error deleting request id '%s'", requestID)
		}
	}

	return nil
}

func SingularityDeplpyExistsHelper(s *terraform.State, client *singularity.Client) error {
	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID
		if _, err := client.GetRequestByID(id); err != nil {
			return fmt.Errorf("Received an error retrieving request id %s", err)
		}
	}
	return nil
}
