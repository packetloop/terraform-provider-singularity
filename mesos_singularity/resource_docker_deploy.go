package mesos_singularity

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	singularity "github.com/lenfree/go-singularity"
)

func resourceDockerDeploy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerDeployCreate,
		Read:   resourceDockerDeployRead,
		Exists: resourceDockerDeployExists,
		Update: resourceDockerDeployUpdate,
		Delete: resourceDockerDeployDelete,

		Schema: map[string]*schema.Schema{
			"deploy_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"args": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"network": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "BRIDGE",
				ValidateFunc: validateDockerNetwork,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
			"force_pull_image": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"num_ports": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  true,
			},
			"envs": envSchema(),
			// We use typeSet because this parameter can be unordered list and must be unique.
			"port_mapping": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"container_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"container_port_type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSingularityPortMappingType,
						},
						"host_port_type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSingularityPortMappingType,
						},
						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSingularityPortProtocol,
							Default:      "tcp",
						},
					},
				},
			},
			"volume": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_path": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"container_path": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"mode": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSingularityDockerVolumeMode,
						},
					},
				},
			},
			"uri": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"cache": &schema.Schema{
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"executable": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"extract": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func envSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func resourceDockerDeployCreate(d *schema.ResourceData, m interface{}) error {

	id := d.Get("deploy_id").(string)
	d.SetId(id)
	return createDockerDeploy(d, m)
}

func resourceDockerDeployExists(d *schema.ResourceData, m interface{}) (b bool, e error) {

	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := clientConn(m)
	r, err := client.GetRequestByID(d.Get("request_id").(string))

	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	if r.RestyResponse.StatusCode() == 404 {
		return false, fmt.Errorf("%v", string(r.RestyResponse.Body()))
	}
	if (r.Body.ActiveDeploy.ID) == "" && (r.Body.RequestDeployState.RequestID == "") {
		return false, fmt.Errorf("%v", string(r.RestyResponse.Body()))
	}
	return true, nil

}

func expandPortMappings(configured []interface{}) ([]singularity.DockerPortMapping, error) {
	var portMappings []singularity.DockerPortMapping
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		l := singularity.DockerPortMapping{
			HostPort:          int64(data["host_port"].(int)),
			ContainerPort:     int64(data["container_port"].(int)),
			ContainerPortType: data["container_port_type"].(string),
			Protocol:          data["protocol"].(string),
			HostPortType:      data["host_port_type"].(string),
		}

		portMappings = append(portMappings, l)
	}
	return portMappings, nil
}

func expandDockerVolumes(configured []interface{}) ([]singularity.SingularityVolume, error) {
	var dockerVolumes []singularity.SingularityVolume
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		l := singularity.SingularityVolume{
			HostPath:      data["host_path"].(string),
			ContainerPath: data["container_path"].(string),
			Mode:          data["mode"].(string),
		}

		dockerVolumes = append(dockerVolumes, l)
	}
	return dockerVolumes, nil
}

func expandUris(configured []interface{}) ([]singularity.SingularityMesosArtifact, error) {
	var uris []singularity.SingularityMesosArtifact
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		l := singularity.SingularityMesosArtifact{
			URI:        data["path"].(string),
			Cache:      data["cache"].(bool),
			Extract:    data["extract"].(bool),
			Executable: data["executable"].(bool),
		}

		uris = append(uris, l)
	}
	return uris, nil
}

func createDockerDeploy(d *schema.ResourceData, m interface{}) error {
	id := strings.ToLower(d.Get("deploy_id").(string))
	requestID := strings.ToLower(d.Get("request_id").(string))
	image := d.Get("image").(string)
	network := d.Get("network").(string)
	cpu := d.Get("cpu").(float64)
	memory := d.Get("memory").(float64)
	numPorts := int64(d.Get("num_ports").(int))
	forcePullImage := d.Get("force_pull_image").(bool)
	command := d.Get("command").(string)
	arguments := d.Get("args").([]interface{})
	env := make(map[string]string)
	envs := d.Get("envs").(map[string]interface{})
	for k, v := range envs {
		env[k] = v.(string)
	}

	portMappings, err := expandPortMappings(d.Get("port_mapping").(*schema.Set).List())

	if int64(len(portMappings)) > numPorts {
		return fmt.Errorf("Error: %s", "Resource num_ports shouldbe >= number of port_mapping")
	}
	d.SetId(id)

	dockerVolumes, err := expandDockerVolumes(d.Get("volume").(*schema.Set).List())
	uris, err := expandUris(d.Get("uri").(*schema.Set).List())

	log.Printf("Singularity deploy '%s' is being provisioned...", id)
	client := clientConn(m)

	info := singularity.ContainerInfo{
		Type: "DOCKER",
		DockerInfo: singularity.DockerInfo{
			ForcePullImage: forcePullImage,
			Network:        strings.ToUpper(network),
			Image:          image,
			PortMappings:   portMappings,
		},
		Volumes: dockerVolumes,
	}
	resource := singularity.SingularityDeployResources{
		Cpus:     cpu,
		MemoryMb: memory,
		NumPorts: numPorts,
	}

	dep := singularity.NewDeploy(id)
	dep.SetURIs(uris)

	// Move this to a map function.
	if len(arguments) > 0 {
		var args []string
		for _, i := range arguments {
			args = append(args, i.(string))
		}
		dep = dep.SetArgs(args...)
	}

	if len(env) > 0 {
		dep = dep.SetEnv(env)
	}

	containerInfo, err := dep.SetContainerInfo(info)

	if err != nil {
		return fmt.Errorf("Create Singularity create deploy error: %v", err)
	}

	deploy := containerInfo.SetCommand(command).
		SetRequestID(requestID).
		SetSkipHealthchecksOnDeploy(true).
		SetResources(resource).
		Build()

	resp, err := singularity.NewDeployRequest().
		AttachDeploy(deploy).
		Build().
		Create(client)

	if err != nil {
		return fmt.Errorf("Create Singularity create deploy error: %v", err)
	}

	return checkDeployResponse(d, m, resp, err)
}

func checkDeployResponse(d *schema.ResourceData, m interface{}, r singularity.HTTPResponse, err error) error {
	log.Printf("[TRACE] check Deploy Response HTTP Response %v", r.RestyResponse)
	if err != nil {
		return fmt.Errorf("Create Singularity deploy error: %v", err)
	}
	if r.RestyResponse.StatusCode() < 200 && r.RestyResponse.StatusCode() > 299 {
		return fmt.Errorf("Create Singularity deploy error %v: %v", r.RestyResponse.StatusCode(), err)
	}
	return resourceDockerDeployRead(d, m)
}

// resourceRequestRead is called to resync the local state with the remote state.
// Terraform guarantees that an existing ID will be set. This ID should be used
// to look up the resource. Any remote data should be updated into the local data.
// No changes to the remote resource are to be made.
func resourceDockerDeployRead(d *schema.ResourceData, m interface{}) error {
	client := clientConn(m)
	r, err := client.GetRequestByID(d.Get("request_id").(string))
	log.Printf("[TRACE] Deploy Read HTTP Response %v", r.Body)

	if err != nil {
		return err
	}
	if r.RestyResponse.StatusCode() == 404 {
		return fmt.Errorf("status code error %v", string(r.RestyResponse.Body()))
	}
	if (r.Body.ActiveDeploy.ID) == "" && (r.Body.RequestDeployState.RequestID == "") {
		return fmt.Errorf("activedeploy empty %v", string(r.RestyResponse.Body()))
	}
	// TODO: Investigate how to deal with Can not deploy a deploy that has already been
	// deployed deploy/task.
	if r.RestyResponse.StatusCode() != 400 && r.Body.RequestType != "SERVICE" {
		d.Set("request_id", r.Body.ActiveDeploy.RequestID)
		d.Set("deploy_id", r.Body.ActiveDeploy.ID)
		d.Set("image", r.Body.ActiveDeploy.Image)
		d.Set("force_pull_image", r.Body.ActiveDeploy.ForcePullImage)
		d.Set("cpu", r.Body.ActiveDeploy.Cpus)
		d.Set("memory", r.Body.ActiveDeploy.MemoryMb)
		d.Set("command", r.Body.ActiveDeploy.Command)
		d.Set("args", r.Body.ActiveDeploy.Arguments)
		d.Set("env", r.Body.ActiveDeploy.Env)
		d.Set("port_mapping", r.Body.ActiveDeploy.PortMappings)
		d.Set("volume", r.Body.ActiveDeploy.Volumes)
		d.Set("uri", r.Body.ActiveDeploy.Uris)
		d.Set("network", strings.ToUpper(r.Body.ActiveDeploy.Network))
	}
	return nil
}

func resourceDockerDeployUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	if d.HasChange("request_id") ||
		d.HasChange("deploy_id") ||
		d.HasChange("image") ||
		d.HasChange("force_pull_image") ||
		d.HasChange("cpu") ||
		d.HasChange("memory") ||
		d.HasChange("args") ||
		d.HasChange("command") ||
		d.HasChange("env") ||
		d.HasChange("port_mapping") ||
		d.HasChange("volume") ||
		d.HasChange("uri") ||
		d.HasChange("network") {
		log.Printf("[TRACE] Create new deploy with request id (%s) success", d.Id())
		d.Partial(false)
		// Singularity deploy is by design idempotent.
		return createDockerDeploy(d, m)
	}
	return nil
}

func resourceDockerDeployDelete(d *schema.ResourceData, m interface{}) error {
	a := deleteRequest(d.Get("request_id").(string))
	d.SetId("")
	return a(d, m)
}

func resourceResourceDockerDeployImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceDockerDeployRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
