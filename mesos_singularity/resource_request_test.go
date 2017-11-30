package mesos_singularity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	singularity "github.com/lenfree/go-mesos-singularity"
)

func TestAccSingularityRequestScheduledCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityRequestScheduledConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_request.foo"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo", "request_id", "foo-test-id"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo", "request_type", "SCHEDULED"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo", "schedule", "0 7 * * *"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo", "schedule_type", "CRON"),
				),
			},
		},
	})
}

func TestAccSingularityRequestWorkerCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityRequestRunOnceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_request.foo-run"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-run", "request_id", "foo-run-id"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-run", "request_type", "RUN_ONCE"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-run", "instances", "5"),
				),
			},
		},
	})
}

const testAccCheckSingularityRequestScheduledConfig = `
resource "singularity_request" "foo" {
			request_id             = "foo-test-id"
			request_type           = "SCHEDULED"
			schedule               = "0 7 * * *"
			schedule_type          = "CRON"
}
`

const testAccCheckSingularityRequestRunOnceConfig = `
resource "singularity_request" "foo-run" {
			request_id             = "foo-run-id"
			request_type           = "RUN_ONCE"
			instances              = 5
}
`

const testAccCheckSingularityRequestServiceConfig = `
resource "singularity_request" "foo-service" {
			request_id             = "foo-service-id"
			request_type           = "SERVICE"
}
`

const testAccCheckSingularityRequestWorkerConfig = `
resource "singularity_request" "foo-worker" {
			request_id             = "foo-worker-id"
			request_type           = "WORKER"
}
`

func testAccCheckSingularityRequestExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Conn).sclient
		return SingularityRequestExistsHelper(s, client)
	}
}

func testCheckSingularityRequestDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "singularity_request" {
			continue
		}

		requestID := res.Primary.ID

		client := testAccProvider.Meta().(*Conn).sclient
		data, err := client.GetRequestByID(requestID)
		if err != nil {
			return nil
		}
		// If request_id does not exists, it gets a response status code 404 Not Found.
		if data.GoRes.StatusCode != 404 {
			return fmt.Errorf("Request id '%s' still exists", requestID)
		}
	}

	return nil
}

func SingularityRequestExistsHelper(s *terraform.State, client *singularity.Client) error {
	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID
		if _, err := client.GetRequestByID(id); err != nil {
			return fmt.Errorf("Received an error retrieving user %s", err)
		}
	}
	return nil
}
