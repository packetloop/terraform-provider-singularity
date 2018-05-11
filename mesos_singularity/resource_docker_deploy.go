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
				ForceNew: true,
				Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "BRIDGE",
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
			"envs": envSchema(),
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
	if (r.Body.ActiveDeploy.ID) == "" {
		return false, fmt.Errorf("%v", string(r.RestyResponse.Body()))
	}
	return true, nil

}

func createDockerDeploy(d *schema.ResourceData, m interface{}) error {
	id := strings.ToLower(d.Get("deploy_id").(string))
	requestID := strings.ToLower(d.Get("request_id").(string))
	image := d.Get("image").(string)
	network := strings.ToUpper(d.Get("network").(string))
	cpu := d.Get("cpu").(float64)
	memory := d.Get("cpu").(float64)
	forcePullImage := d.Get("force_pull_image").(bool)
	command := d.Get("command").(string)
	args2 := d.Get("args").([]interface{})
	var args []string
	for _, i := range args2 {
		args = append(args, i.(string))
	}

	env := make(map[string]string)
	envs := d.Get("envs").(map[string]interface{})
	for k, v := range envs {
		env[k] = v.(string)
	}

	d.SetId(id)

	log.Printf("Singularity deploy '%s' is being provisioned...", id)
	client := clientConn(m)

	info := singularity.ContainerInfo{
		Type: "DOCKER",
		DockerInfo: singularity.DockerInfo{
			ForcePullImage: forcePullImage,
			Network:        network,
			Image:          image,
		},
	}
	resource := singularity.SingularityDeployResources{
		Cpus:     cpu,
		MemoryMb: memory,
	}
	dep := singularity.NewDeploy(id)
	containerInfo, err := dep.SetContainerInfo(info)
	if err != nil {
		return fmt.Errorf("Create Singularity create deploy error: %v", err)
	}
	deploy := containerInfo.SetCommand(command).
		SetArgs(args...).
		SetRequestID(requestID).
		SetSkipHealthchecksOnDeploy(true).
		SetEnv(env).
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
	log.Printf("[TRACE] HTTP Response %v", r.RestyResponse)

	if err != nil {
		return fmt.Errorf("Create Singularity deploy error: %v", err)
	}
	if r.RestyResponse.StatusCode() <= 200 && r.RestyResponse.StatusCode() >= 299 {
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

	if err != nil {
		return err
	}
	if r.RestyResponse.StatusCode() == 404 {
		return fmt.Errorf("%v", string(r.RestyResponse.Body()))
	}
	if (r.Body.ActiveDeploy.ID) == "" {
		return fmt.Errorf("%v", string(r.RestyResponse.Body()))
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
		d.HasChange("network") {
		log.Printf("[TRACE] Delete and update existing request id (%s) success", d.Id())
		// TODO: Investigate whether we can just update existing request, rather
		// than delete and add.
		// Delete existing request and add if there are changes. I couldn't manage
		// to find API doco to update existing request. Only for existing deploy.

		d.Partial(false)
		return resourceDockerDeployCreate(d, m)
	}
	return nil
}

func resourceDockerDeployDelete(d *schema.ResourceData, m interface{}) error {
	a := deleteRequest(d.Get("request_id").(string))
	d.SetId("")
	return a(d, m)
}
