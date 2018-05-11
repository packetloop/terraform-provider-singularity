package mesos_singularity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	singularity "github.com/lenfree/go-singularity"
)

func TestAccSingularityDockerDeployCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityDeployDockerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_.foo"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "deploy_id", "mydeploy"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "force_pull_image", "false"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "network", "bridge"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "image", "golang:latest"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "cpu", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "memory", "128"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "command", "bash"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "request_id", "myrequest"),
				),
			},
		},
	})
}

const testAccCheckSingularityDeployDockerConfig = `
resource "singularity_request" "foo" {
	request_id             = "myrequest"
	request_type           = "SCHEDULED"
	schedule               = "0 7 * * *"
	schedule_type          = "CRON"
}
resource "singularity_docker_deploy" "foo" {
			deploy_id        = "mydeploy"
			force_pull_image = false
			network          = "bridge"
			image            = "golang:latest"
			cpu              = 2
			memory           = 128
			command          = "bash"
			args             = ["-xc", "date"]
			request_id       = "${singularity_request.foo.id}"
}
`

func testAccCheckSingularityDockerDeployExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Conn).sclient
		return SingularityDockerDeployExistsHelper(s, client)
	}
}

func testCheckSingularityDockerDeployDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "singularity_request" {
			continue
		}

		requestID := res.Primary.Attributes["request_id"]

		client := testAccProvider.Meta().(*Conn).sclient
		data, err := client.GetRequestByID(requestID)
		if err != nil {
			return err
		}
		// If request_id does not exists, it gets a response status code 404 Not Found.
		if data.RestyResponse.StatusCode() != 404 {
			continue
		}
		return fmt.Errorf("Request id '%s' still exists", requestID)
	}
	return nil
}

func SingularityDockerDeployExistsHelper(s *terraform.State, client *singularity.Client) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "singularity_docker_deploy" {
			continue
		}
		reqID := res.Primary.Attributes["request_id"]
		//id := r.Primary.ID
		if _, err := client.GetRequestByID(reqID); err != nil {
			return fmt.Errorf("Received an error retrieving request id %v", err)
		}
	}
	return nil
}
