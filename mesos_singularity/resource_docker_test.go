package mesos_singularity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	singularity "github.com/lenfree/go-singularity"
)

func TestAccSingularityDeployDockerCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityDeployConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityDeployDockerExists("singularity_deploy_docker.foo"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo", "request_id", "foo-test-id"),
				),
			},
		},
	})
}

const testAccCheckSingularityDeployDockerConfig = `
resource "singularity_deploy_docker" "foo" {
  	image = golang:latest
  	network = BRIDGE
  	portMappings = [{
  	  containerPortType = LITERAL
  	  containerPort = 3000
  	  hostPortType = FROM_OFFER
  	  hostPort = 0
	  protocol = tcp
	]
}
`

func testAccCheckSingularityDeployDockerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Conn).sclient
		return SingularityRequestExistsHelper(s, client)
	}
}

func testCheckSingularityDeployDockerDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "singularity_deploy_docker" {
			continue
		}

		requestID := res.Primary.ID

		client := testAccProvider.Meta().(*Conn).sclient
		data, err := client.GetRequestByID(requestID)
		if err != nil {
			return nil
		}
		// If request_id does not exists, it gets a response status code 404 Not Found.
		if data.RestyResponse.StatusCode() != 404 {
			return fmt.Errorf("Request id '%s' still exists", requestID)
		}
	}

	return nil
}

func SingularityDeployDockerExistsHelper(s *terraform.State, client *singularity.Client) error {
	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID
		if _, err := client.GetRequestByID(id); err != nil {
			return fmt.Errorf("Received an error retrieving request id %s", err)
		}
	}
	return nil
}
