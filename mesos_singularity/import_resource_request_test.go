package mesos_singularity

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSingularityRequest_importRequestType(t *testing.T) {
	resourceName := "singularity_request.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSingularityRequestScheduledConfig,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// We ignore num_retries_on_failure because not all request type
				// accepts this parameter.
				ImportStateVerifyIgnore: []string{
					"num_retries_on_failure",
				},
			},
		},
	})
}
