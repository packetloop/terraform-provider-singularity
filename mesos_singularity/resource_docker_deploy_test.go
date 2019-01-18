package mesos_singularity

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
						"singularity_docker_deploy.foo", "container_info.0.docker_info.#", "1"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "container_info.0.docker_info.0.image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "container_info.0.docker_info.0.network", "BRIDGE"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "resources.cpus", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foo", "resources.memory_mb", "128"),
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
						"singularity_docker_deploy.bar", "container_info.0.docker_info.#", "1"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "container_info.0.docker_info.0.image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "container_info.0.docker_info.0.network", "BRIDGE"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "resources.cpus", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.bar", "resources.memory_mb", "128"),
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
						"singularity_docker_deploy.foobar", "deploy_id", "mydeployfoobar330c"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "container_info.0.docker_info.#", "1"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "container_info.0.docker_info.0.image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "container_info.0.docker_info.0.network", "BRIDGE"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "resources.cpus", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "resources.memory_mb", "128"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "command", "bash"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "request_id", "myrequestfoobar"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobar", "container_info.0.docker_info.0.port_mapping.#", "2"),
					// Skip attribute test for port_mapping because schema.typeSet items are
					// are stored in state with an index value calculated by the hash of the
					// attributes of the set according to
					// https://www.terraform.io/docs/extend/schemas/schema-types.html
				),
			},
		},
	})
}

func TestAccSingularityDockerDeployCreateVolumes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckSingularityRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSingularityDeployDockerConfigVolumes,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSingularityRequestExists("singularity_deploy.foobaz"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "deploy_id", "mydeployfoobazz2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "container_info.0.docker_info.#", "1"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "container_info.0.docker_info.0.image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "container_info.0.docker_info.0.network", "BRIDGE"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "resources.cpus", "2"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "resources.memory_mb", "128"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "command", "bash"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "request_id", "myrequestfoobaz"),
					resource.TestCheckResourceAttr(
						"singularity_docker_deploy.foobaz", "container_info.0.volume.#", "2"),
					// Skip attribute test for volume because schema.typeSet items are
					// are stored in state with an index value calculated by the hash of the
					// attributes of the set according to
					// https://www.terraform.io/docs/extend/schemas/schema-types.html
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
  command          = "bash"
  args             = ["-xc", "date"]
  request_id       = "${singularity_request.foo.id}"

  container_info {
    docker_info {
      force_pull_image = false
      network          = "BRIDGE"
      image            = "ubuntu"
    }
  }
  resources {
    cpus      = 2
    memory_mb = 128
  }
}
`

const testAccCheckSingularityDeployDockerConfigMaxTasks = `
resource "singularity_request" "bar" {
  request_id             = "myrequestbar"
  request_type           = "SCHEDULED"
  schedule               = "0 7 * * *"
  schedule_type          = "CRON"
}
resource "singularity_docker_deploy" "bar" {
  deploy_id        = "mydeploybar"
  command          = "bash"
  args             = ["-xc", "date"]
  request_id       = "${singularity_request.bar.id}"

  container_info {
    docker_info {
      force_pull_image = false
      network          = "BRIDGE"
      image            = "ubuntu"
    }
  }
  envs {
    MYENV = "test"
    NAME  = "lenfree"
  }

  resources {
    cpus      = 2
    memory_mb = 128
  }
}
`

const testAccCheckSingularityDeployDockerConfigPortMapping = `
resource "singularity_request" "foobar" {
  request_id             = "myrequestfoobar"
  request_type           = "SERVICE"
  instances              = 2
}
resource "singularity_docker_deploy" "foobar" {
  deploy_id        = "mydeployfoobar330c"
  command          = "bash"
  args             = ["-xc", "while true; do echo up; done"]
  request_id       = "${singularity_request.foobar.id}"

  container_info {
    docker_info {
      force_pull_image = false
      network          = "BRIDGE"
      image            = "ubuntu"

	  port_mapping {
    	host_port           = 1
    	container_port      = 9999
    	container_port_type = "LITERAL"
    	host_port_type      = "FROM_OFFER"
    	protocol            = "tcp"
      }

      port_mapping {
    	host_port           = 0
    	container_port      = 10001
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

const testAccCheckSingularityDeployDockerConfigVolumes = `
resource "singularity_request" "foobaz" {
  request_id             = "myrequestfoobaz"
  request_type           = "SERVICE"
  instances              = 1
}
resource "singularity_docker_deploy" "foobaz" {
  deploy_id        = "mydeployfoobazz2"
  args             = ["-xc", "while true; do echo up; done"]
  request_id       = "${singularity_request.foobaz.id}"
  command          = "bash"

  container_info {
    docker_info {
      force_pull_image = "false"
      network          = "BRIDGE"
      image            = "ubuntu"
    }
    volume {
      mode           = "RO"
      container_path = "/root/.aws/config"
      host_path      = "/root/.aws/config"
    }
    volume {
      mode           = "RO"
      container_path = "/root/.aws/credentials"
      host_path      = "/root/.aws/credentials"
    }
  }

  resources {
    cpus      = 2
    memory_mb = 128
  }
}
`

func testAccCheckSingularityDockerDeployExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Conn).sclient
		return SingularityDockerDeployExistsHelper(s, client)
	}
}
func testCheckSingularityDockerDeployDestroyTest(state *terraform.State) error {
	return nil
}

func testCheckSingularityDockerDeployDestroy(state *terraform.State) error {
	//for _, res := range state.RootModule().Resources {
	//	if res.Type != "singularity_destroy" {
	//		continue
	//	}

	//	requestID := res.Primary.Attributes["request_id"]

	//	client := testAccProvider.Meta().(*Conn).sclient
	//	data, err := client.GetRequestByID(requestID)
	//	if err != nil {
	//		return err
	//	}
	//	// If request_id does not exists, it gets a response status code 404 Not Found.
	//	if data.RestyResponse.StatusCode() != 404 {
	//		continue
	//	}
	//	return fmt.Errorf("Request id '%s' still exists", requestID)
	//}
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

func TestExpandPortMappings(t *testing.T) {
	portMappings := []struct {
		val    []interface{}
		expect []singularity.DockerPortMapping
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"host_port":           0,
					"container_port":      10001,
					"container_port_type": "LITERAL",
					"host_port_type":      "FROM_OFFER",
					"protocol":            "tcp",
				},
			},
			[]singularity.DockerPortMapping{
				singularity.DockerPortMapping{
					HostPort:          0,
					ContainerPort:     10001,
					ContainerPortType: "LITERAL",
					HostPortType:      "FROM_OFFER",
					Protocol:          "tcp",
				},
			},
		},
	}

	for _, data := range portMappings {
		s := schema.NewSet(portMappingHash, []interface{}{})
		for _, v := range data.val {
			s.Add(v)
		}
		actual := expandPortMappings(s)
		if diff := reflect.DeepEqual(data.expect, actual); !diff {
			t.Errorf("Got %+v\n, wants %#+v\n, actual %#+v\n, passed %v\n", diff, data.expect, actual, data.val)
		}
	}
}

func TestExpandDockerVolumes(t *testing.T) {
	volumes := []struct {
		val    []interface{}
		expect []singularity.SingularityVolume
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"mode":           "RW",
					"container_path": "/inside/path",
					"host_path":      "/outside/path",
				},
			},
			[]singularity.SingularityVolume{
				singularity.SingularityVolume{
					Mode:          "RW",
					ContainerPath: "/inside/path",
					HostPath:      "/outside/path",
				},
			},
		},
	}
	for _, data := range volumes {
		s := schema.NewSet(containerVolumeHash, []interface{}{})
		for _, v := range data.val {
			s.Add(v)
		}
		actual := expandContainerVolumes(s)
		if diff := reflect.DeepEqual(data.expect, actual); !diff {
			t.Errorf("Got %+v\n, wants %#+v\n, actual %#+v\n, passed %v\n", diff, data.expect, actual, data.val)
		}
	}
}

func TestExpandUris(t *testing.T) {
	portMappings := []struct {
		val    []interface{}
		expect []singularity.SingularityMesosArtifact
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"path":       "file:///etc/docker.tar.gz",
					"cache":      false,
					"executable": false,
					"extract":    true,
				},
			},
			[]singularity.SingularityMesosArtifact{
				singularity.SingularityMesosArtifact{
					URI:        "file:///etc/docker.tar.gz",
					Cache:      false,
					Executable: false,
					Extract:    true,
				},
			},
		},
	}
	for _, data := range portMappings {
		actual, err := expandUris(data.val)
		if err != nil {
			t.Errorf("Error %v\n", err)
		}
		if diff := reflect.DeepEqual(data.expect, actual); !diff {
			t.Errorf("Got %+v\n, wants %#+v\n, actual %#+v\n, passed %v\n", diff, data.expect, actual, data.val)
		}
	}
}
