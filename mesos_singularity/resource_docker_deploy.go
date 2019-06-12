package mesos_singularity

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"crypto/md5"
	"github.com/cydev/zero"
	"github.com/hashicorp/terraform/helper/hashcode"
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
		Importer: &schema.ResourceImporter{
			State: resourceResourceDockerDeployImport,
		},

		Schema: map[string]*schema.Schema{
			"deploy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"args": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"resources": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory_mb": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"cpus": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			"container_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"docker_info": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"force_pull_image": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  "false",
									},
									"network": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "BRIDGE",
										ValidateFunc: validateDockerNetwork,
									},
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
					},
				},
			},
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
			"envs":     envSchema(),
			"metadata": envSchema(),
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
	if r.RestyResponse.StatusCode() == 400 {
		return false, fmt.Errorf("%v", string(r.RestyResponse.Body()))
	}
	if (r.Body.ActiveDeploy.ID) == "" && (r.Body.RequestDeployState.RequestID == "") {
		return false, fmt.Errorf("%v", string(r.RestyResponse.Body()))
	}
	return true, nil
}

func expandContainerVolume(v map[string]interface{}) singularity.SingularityVolume {
	return singularity.SingularityVolume{
		HostPath:      v["host_path"].(string),
		ContainerPath: v["container_path"].(string),
		Mode:          v["mode"].(string),
	}
}

func expandContainerVolumes(configured *schema.Set) []singularity.SingularityVolume {
	c := configured.List()
	var dockerVolumes []singularity.SingularityVolume
	for _, lRaw := range c {
		data := lRaw.(map[string]interface{})
		dockerVolumes = append(dockerVolumes, expandContainerVolume(data))
	}
	return dockerVolumes
}
func expandVolumes(d map[string]interface{}) []singularity.SingularityVolume {
	v := d["volume"].(*schema.Set)
	return expandContainerVolumes(v)
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

func expandResources(d *schema.ResourceData, portMappings int64) (singularity.SingularityDeployResources, error) {

	cpus, err := strconv.ParseFloat(d.Get("resources.cpus").(string), 64)
	if err != nil {
		return singularity.SingularityDeployResources{}, fmt.Errorf("Error converting cpus to float64: %v", err)
	}
	memoryMb, err := strconv.ParseFloat(d.Get("resources.memory_mb").(string), 64)
	if err != nil {
		return singularity.SingularityDeployResources{}, fmt.Errorf("Error converting memory_mb to float64: %v", err)
	}

	return singularity.SingularityDeployResources{
		Cpus:     cpus,
		MemoryMb: memoryMb,
		NumPorts: 2,
	}, nil
}

func expandPortMappings(configured *schema.Set) []singularity.DockerPortMapping {
	p := configured.List()
	var portMappings []singularity.DockerPortMapping

	for _, mRaw := range p {
		data := mRaw.(map[string]interface{})

		l := singularity.DockerPortMapping{
			HostPort:          data["host_port"].(int),
			ContainerPort:     data["container_port"].(int),
			ContainerPortType: data["container_port_type"].(string),
			Protocol:          data["protocol"].(string),
			HostPortType:      data["host_port_type"].(string),
		}

		portMappings = append(portMappings, l)
	}
	return portMappings
}

func expandDockerInfo(d map[string]interface{}) singularity.DockerInfo {
	a := d["docker_info"].([]interface{})
	var portMappings []singularity.DockerPortMapping
	var forcePullImage bool
	var network string
	var image string
	for _, i := range a {
		if i, ok := i.(map[string]interface{}); ok {
			forcePullImage = i["force_pull_image"].(bool)
			network = i["network"].(string)
			image = i["image"].(string)
			pm := i["port_mapping"].(*schema.Set)
			portMappings = expandPortMappings(pm)
		}
	}
	return singularity.DockerInfo{
		ForcePullImage: forcePullImage,
		Network:        network,
		Image:          image,
		PortMappings:   portMappings,
	}
}

func expandContainerInfo(d *schema.ResourceData) singularity.ContainerInfo {
	a := d.Get("container_info").([]interface{})

	var dockerInfo singularity.DockerInfo
	var volumes []singularity.SingularityVolume
	for _, i := range a {
		dockerInfo = expandDockerInfo(i.(map[string]interface{}))
		volumes = expandVolumes(i.(map[string]interface{}))
	}
	return singularity.ContainerInfo{
		Type:       "DOCKER",
		DockerInfo: dockerInfo,
		Volumes:    volumes,
	}
}

func resourceDockerDeployCreate(d *schema.ResourceData, m interface{}) error {

	id := calculateDeployMD5(buildDeployRequest(d))
	d.SetId(id)
	return createDockerDeploy(d, m)
}

func calculateDeployMD5(input singularity.DeployRequest) string{
	data := []byte(fmt.Sprintf("%v", input))
	return fmt.Sprintf("%x", md5.Sum(data))
}

func buildDeployRequest(d *schema.ResourceData) singularity.DeployRequest {
	requestID := strings.ToLower(d.Get("request_id").(string))
	command := d.Get("command").(string)
	arguments := d.Get("args").([]interface{})
	env := make(map[string]string)
	envs := d.Get("envs").(map[string]interface{})
	for k, v := range envs {
		env[k] = v.(string)
	}

	dockerVolumes, err := expandDockerVolumes(d.Get("volume").(*schema.Set).List())
	uris, _:= expandUris(d.Get("uri").(*schema.Set).List())

	info := expandContainerInfo(d)

	resources, _:= expandResources(d, int64(len(info.DockerInfo.PortMappings)))

	dep := singularity.NewDeploy("")
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

	containerInfo, _ := dep.SetContainerInfo(info)

	deploy := containerInfo.SetCommand(command).
		SetRequestID(requestID).
		SetResources(resources).
		SetSkipHealthchecksOnDeploy(true).
		Build()

	resp := singularity.NewDeployRequest().
		AttachDeploy(deploy).
		Build()

	return resp
}

func createDockerDeploy(d *schema.ResourceData, m interface{}) error {
	client := clientConn(m)
	// Workaround update ID with md5sum of config params
	md5 := calculateDeployMD5(buildDeployRequest(d))
	deployRequest := buildDeployRequest(d).SetID(md5)

	log.Printf("Singularity deploy '%s' is being provisioned...", md5)
	resp, err := deployRequest.Create(client)
	if err != nil {
		return fmt.Errorf("Singularity ID: %v deploy error: %v", d.Get("ID"), err)
	}


	return checkDeployResponse(d, m, resp, err)
}

func checkDeployResponse(d *schema.ResourceData, m interface{}, r singularity.HTTPResponse, err error) error {
	log.Printf("[TRACE] check Deploy Response HTTP Response %v", r.RestyResponse)
	if err != nil {
		return fmt.Errorf("Create Singularity deploy error: %v", err)
	}

	if r.RestyResponse.StatusCode() < 200 && r.RestyResponse.StatusCode() > 299 {
		return fmt.Errorf("Create Singularity deploy error: %v, %v", r.RestyResponse.StatusCode(), err)
	}
	return resourceDockerDeployRead(d, m)
}

// resourceRequestRead is called to resync the local state with the remote state.
// Terraform guarantees that an existing ID will be set. This ID should be used
// to look up the resource. Any remote data should be updated into the local data.
// No changes to the remote resource are to be made.
func resourceDockerDeployRead(d *schema.ResourceData, m interface{}) error {
	client := clientConn(m)
	//deploy_id := d.Get("deploy_id").(string)
	//r, err := client.GetRequestByID(d.Get("request_id").(string))
	//log.Printf("[TRACE] Deploy Read HTTP Response %v", r.Body)

	// Expensive loop. Only use this during import because we don't have access to other attributes than
	// GetID(). Otherwise, use getrequestsbyid.
	_, b, err := client.GetRequests()
	if err != nil {
		d.SetId("")
		return err
	}
	log.Printf("[TRACE] GETT %v", d)
	id := d.Id()
	c := b.GetRequestID(id)
	//log.Printf("[TRACE] Deploy Read HTTP Response %v", string(res.Body()))
	r, err := client.GetRequestByID(c.SingularityRequest.ID)
	if err != nil {
		d.SetId("")
		return err
	}

	// When we create a service request, a deploy does not run immediately by default
	// and deploy would be in pending state. Hence, we just check if struct is empty
	// and if it is empty, we use activedeploy object instead.
	if zero.IsZero(r.Body.ActiveDeploy.ID) {
		d.Set("deploy_id", r.Body.RequestDeployState.PendingDeployState.DeployID)
		d.Set("command", r.Body.PendingDeploy.Command)
		d.Set("envs", r.Body.PendingDeploy.TaskEnv)

		cpus := strconv.FormatFloat(r.Body.PendingDeploy.Cpus, 'f', -1, 64)
		memoryMb := strconv.FormatFloat(r.Body.PendingDeploy.MemoryMb, 'f', -1, 64)

		resources := make(map[string]string)
		for k, v := range map[string]string{
			"cpus":      cpus,
			"memory_mb": memoryMb,
		} {
			resources[k] = v
		}
		d.Set("resources", resources)

		if r.Body.PendingDeploy.Uris != nil {
			mapURI := make([]map[string]interface{}, 0)
			for _, a := range r.Body.PendingDeploy.Uris {
				m := make(map[string]interface{})
				m["cache"] = a.Cache
				m["path"] = a.URI
				m["extract"] = a.Extract
				m["executable"] = a.Executable
				mapURI = append(mapURI, m)
			}
		}
		d.Set("uri", r.Body.PendingDeploy.Uris)
		d.Set("metadata", r.Body.PendingDeploy.Metadata)
		if err = d.Set("container_info", flattenContainerInfo(r.Body.PendingDeploy.ContainerInfo)); err != nil {
			return fmt.Errorf("flatten docker_info from pendingDeploy error: %v", err)
		}
	} else {
		d.Set("deploy_id", r.Body.ActiveDeploy.ID)
		d.Set("args", r.Body.ActiveDeploy.Arguments)
		d.Set("command", r.Body.ActiveDeploy.Command)
		d.Set("envs", r.Body.ActiveDeploy.Env)

		cpus := strconv.FormatFloat(r.Body.ActiveDeploy.Cpus, 'f', -1, 64)
		memoryMb := strconv.FormatFloat(r.Body.ActiveDeploy.MemoryMb, 'f', -1, 64)

		resources := make(map[string]string)
		for k, v := range map[string]string{
			"cpus":      cpus,
			"memory_mb": memoryMb,
		} {
			resources[k] = v
		}
		d.Set("resources", resources)

		if r.Body.ActiveDeploy.Uris != nil {
			mapURI := make([]map[string]interface{}, 0)
			for _, a := range r.Body.ActiveDeploy.Uris {
				m := make(map[string]interface{})
				m["cache"] = a.Cache
				m["path"] = a.URI
				m["extract"] = a.Extract
				m["executable"] = a.Executable
				mapURI = append(mapURI, m)
			}
			d.Set("uri", mapURI)
		}
		d.Set("metadata", r.Body.ActiveDeploy.Metadata)

		if err = d.Set("container_info", flattenContainerInfo(r.Body.ActiveDeploy.ContainerInfo)); err != nil {
			return fmt.Errorf("flatten docker_info from activeDeploy error: %v", err)
		}
		d.Set("args", r.Body.ActiveDeploy.Arguments)
	}
	d.Set("request_id", r.Body.SingularityRequest.ID)
	return nil
}

func flattenContainerInfo(in singularity.ContainerInfo) []interface{} {
	m := make(map[string]interface{})
	m["docker_info"] = flattenDockerInfo(in.DockerInfo)
	m["volume"] = flattenContainerVolumes(in.Volumes)
	return []interface{}{m}
}

func flattenDockerInfo(in singularity.DockerInfo) []interface{} {
	m := make(map[string]interface{})
	m["network"] = in.Network
	m["image"] = in.Image
	m["force_pull_image"] = in.ForcePullImage
	m["port_mapping"] = flattenDockerPortMappings(in.PortMappings)
	return []interface{}{m}
}
func flattenContainerVolumes(in []singularity.SingularityVolume) *schema.Set {
	s := schema.NewSet(containerVolumeHash, []interface{}{})
	for _, v := range in {
		s.Add(flattenContainerVolume(v))
	}
	return s
}
func containerVolumeHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["host_path"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["container_path"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mode"].(string)))
	return hashcode.String(buf.String())
}

func flattenContainerVolume(v singularity.SingularityVolume) map[string]interface{} {
	m := make(map[string]interface{})
	m["host_path"] = v.HostPath
	m["container_path"] = v.ContainerPath
	m["mode"] = v.Mode
	return m
}

func portMappingHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["container_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", m["container_port_type"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m["host_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", m["host_port_type"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["protocol"].(string)))
	return hashcode.String(buf.String())
}

func flattenDockerPortMappings(in []singularity.DockerPortMapping) *schema.Set {
	s := schema.NewSet(portMappingHash, []interface{}{})
	for _, v := range in {
		s.Add(flattenDockerPortMapping(v))
	}
	return s
}

func flattenDockerPortMapping(v singularity.DockerPortMapping) map[string]interface{} {
	m := make(map[string]interface{})
	m["container_port"] = v.ContainerPort
	m["container_port_type"] = v.ContainerPortType
	m["host_port"] = v.HostPort
	m["host_port_type"] = v.HostPortType
	m["protocol"] = v.Protocol
	return m
}

func resourceDockerDeployUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	if d.HasChange("request_id") ||
		d.HasChange("deploy_id") ||
		d.HasChange("container_info") ||
		d.HasChange("resources") ||
		d.HasChange("args") ||
		d.HasChange("command") ||
		d.HasChange("envs") ||
		d.HasChange("metadata") ||
		d.HasChange("uri") {
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
