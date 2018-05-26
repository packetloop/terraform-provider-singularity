package mesos_singularity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	singularity "github.com/lenfree/go-singularity"
)

func TestAccSingularityDockerDeployCreateDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityDeployDockerConfigDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_deploy.foo"),
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
func TestAccSingularityDockerDeployCreateMaxOffer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityDeployDockerConfigMaxTasks,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_deploy.bar"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "deploy_id", "mydeploybar"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "force_pull_image", "false"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "network", "bridge"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "image", "golang:latest"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "cpu", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "memory", "128"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "command", "bash"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "request_id", "myrequestbar"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "envs.%", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "envs.MYENV", "test"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "envs.NAME", "lenfree"),
				),
			},
		},
	})
}

func TestAccSingularityDockerDeployCreatePortMapping(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityDeployDockerConfigPortMapping,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_deploy.foobar"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "deploy_id", "mydeployfoobar"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "force_pull_image", "false"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "network", "bridge"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "image", "golang:latest"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "cpu", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "memory", "128"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "num_ports", "1"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "command", "bash"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "request_id", "myrequestfoobar"),
				),
			},
		},
	})
}

const testAccCheckSingularityDeployDockerConfigDefault = `
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
			envs {
				MYENV = "test"
				NAME  = "lenfree"
			}
}
`

const testAccCheckSingularityDeployDockerConfigMaxTasks = `
resource "singularity_request" "bar" {
	request_id             = "myrequestbar"
	request_type           = "SCHEDULED"
	schedule               = "0 7 * * *"
	schedule_type          = "CRON"
	max_tasks_per_offer    = 2
}
resource "singularity_docker_deploy" "bar" {
	deploy_id        = "mydeploybar"
	force_pull_image = false
	network          = "bridge"
	image            = "golang:latest"
	cpu              = 2
	memory           = 128
	command          = "bash"
	args             = ["-xc", "date"]
	request_id       = "${singularity_request.bar.id}"
	envs {
		MYENV = "test"
		NAME  = "lenfree"
	}
}
`

const testAccCheckSingularityDeployDockerConfigPortMapping = `
resource "singularity_request" "foobar" {
	request_id             = "myrequestfoobar"
	request_type           = "SERVICE"
	instances              = 1
	max_tasks_per_offer    = 2
}
resource "singularity_docker_deploy" "foobar" {
	deploy_id        = "mydeployfoobar"
	force_pull_image = false
	network          = "bridge"
	image            = "golang:latest"
	cpu              = 2
	memory           = 128
	num_ports        = 1
	command          = "bash"
	args             = ["-xc", "date"]
	request_id       = "${singularity_request.foobar.id}"
	envs {
		"MYENV" = "test"
		"NAME"  = "lenfree"
	}
	port_mapping {
		host_port           = 0
		container_port      = 10001
		container_port_type = "LITERAL"
		host_port_type      = "FROM_OFFER"
		protocol            = "tcp"
	}
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
