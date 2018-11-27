package mesos_singularity

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"singularity": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("HOST"); v == "" {
		log.Println("[INFO] Test: Using 'localhost' as test host")
		os.Setenv("HOST", "localhost")
	}
	if v := os.Getenv("PORT"); v != "" {
		log.Printf("[INFO] Test: Using port '%s'", os.Getenv("PORT"))
		os.Setenv("PORT", os.Getenv("PORT"))
	}
	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}
