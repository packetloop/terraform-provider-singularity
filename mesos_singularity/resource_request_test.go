package mesos_singularity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	singularity "github.com/lenfree/go-singularity"
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

func TestAccSingularityRequestRunOnceCreate(t *testing.T) {
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
func TestAccSingularityRequestServiceCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityRequestServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_request.bar"),
					resource.TestCheckResourceAttr(
						"singularity_request.bar", "request_id", "foo-service-id"),
					resource.TestCheckResourceAttr(
						"singularity_request.bar", "request_type", "SERVICE"),
					resource.TestCheckResourceAttr(
						"singularity_request.bar", "instances", "3"),
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
				Config: testAccCheckSingularityRequestWorkerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_request.foo-worker"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-worker", "request_id", "foo-worker-id"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-worker", "request_type", "WORKER"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-worker", "instances", "2"),
				),
			},
		},
	})
}

func TestAccSingularityRequesOnDemandCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityRequestOnDemandConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_request.foo-ondemand"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-ondemand", "request_id", "foo-ondemand-id"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-ondemand", "request_type", "ON_DEMAND"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-ondemand", "instances", "2"),
				),
			},
		},
	})
}

func TestAccSingularityRequesOnDemandCreateDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityRequestOnDemandConfigDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_request.foo-default"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-default", "request_id", "foo-default-id"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-default", "request_type", "ON_DEMAND"),
					resource.TestCheckResourceAttr(
						"singularity_request.foo-default", "instances", "1"),
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
resource "singularity_request" "bar" {
			request_id             = "foo-service-id"
			request_type           = "SERVICE"
			instances              = 3
}
`

const testAccCheckSingularityRequestWorkerConfig = `
resource "singularity_request" "foo-worker" {
			request_id             = "foo-worker-id"
			request_type           = "WORKER"
			instances              = 2
}
`

const testAccCheckSingularityRequestOnDemandConfig = `
resource "singularity_request" "foo-ondemand" {
			request_id             = "foo-ondemand-id"
			request_type           = "ON_DEMAND"
			instances              = 2
}
`

const testAccCheckSingularityRequestOnDemandConfigDefault = `
resource "singularity_request" "foo-default" {
			request_id             = "foo-default-id"
			request_type           = "ON_DEMAND"
			instances              = 1
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
			return err
		}
		// If request_id does not exists, it gets a response status code 404 Not Found.
		if data.RestyResponse.StatusCode() != 404 {
			continue
		}
		return fmt.Errorf("Request id '%s' still exists, %v", requestID, data.RestyResponse.StatusCode())
	}

	return nil
}

func SingularityRequestExistsHelper(s *terraform.State, client *singularity.Client) error {
	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID
		if _, err := client.GetRequestByID(id); err != nil {
			return fmt.Errorf("Received an error retrieving request id %s", err)
		}
	}
	return nil
}

const testAccCheckSingularityRequestSlavePlacement = `
resource "singularity_request" "fooslaveplacement" {
			request_id             = "fooslaveplacement-test-id"
			request_type           = "SCHEDULED"
			schedule               = "0 7 * * *"
			schedule_type          = "CRON"
			slave_placement        = "SEPARATE_BY_DEPLOY"
			num_retries_on_failure = 3
}
`

func TestAccSingularityRequestSlavePlacement(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityRequestSlavePlacement,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_request.fooslaveplacement"),
					resource.TestCheckResourceAttr(
						"singularity_request.fooslaveplacement", "request_id", "fooslaveplacement-test-id"),
					resource.TestCheckResourceAttr(
						"singularity_request.fooslaveplacement", "request_type", "SCHEDULED"),
					resource.TestCheckResourceAttr(
						"singularity_request.fooslaveplacement", "schedule", "0 7 * * *"),
					resource.TestCheckResourceAttr(
						"singularity_request.fooslaveplacement", "schedule_type", "CRON"),
					resource.TestCheckResourceAttr(
						"singularity_request.fooslaveplacement", "slave_placement", "SEPARATE_BY_DEPLOY"),
				),
			},
		},
	})
}
