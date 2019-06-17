package mesos_singularity

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	singularity "github.com/lenfree/go-singularity"
)

func resourceRequest() *schema.Resource {
	return &schema.Resource{
		Create: resourceRequestCreate,
		Read:   resourceRequestRead,
		Exists: resourceRequestExists,
		Update: resourceRequestUpdate,
		Delete: resourceRequestDelete,
		Importer: &schema.ResourceImporter{
			State: resourceResourceRequestImport,
		},

		Schema: map[string]*schema.Schema{
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRequestType,
			},
			"num_retries_on_failure": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"schedule_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRequestScheduleType,
				ForceNew:     true,
			},
			"instances": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"state": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRequestState,
			},
			"slave_placement": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SEPARATE_BY_DEPLOY",
				ForceNew:     true,
				ValidateFunc: validateRequestSlavePlacement,
			},
		},
	}
}

func resourceRequestCreate(d *schema.ResourceData, m interface{}) error {
	id := d.Get("request_id").(string)
	d.SetId(id)
	log.Printf("[INFO] Creating request id: (%s)", id)
	return createRequest(d, m)
}

func resourceScaleRequest(d *schema.ResourceData, m interface{}) error {
	client := clientConn(m)
	id := d.Get("request_id").(string)
	instances := d.Get("instances").(int)
	message := fmt.Sprintf("scale to %d", instances)
	// TODO:
	// Make this configurable
	increment := 1
	log.Printf("[INFO] Scale request id: (%s)", id)
	req := singularity.NewRequestScale(
		id,
		message,
		instances,
		increment,
	)
	resp, err := singularity.ScaleRequest(client, *req)
	if err != nil {
		return fmt.Errorf("scale request ID: (%v) error, %v", id, err)
	}
	if resp.RestyResponse.StatusCode() == 200 {
		return nil
	}
	return fmt.Errorf("scale request ID: (%v) error, %v", id, err)
}

func resourceRequestExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := clientConn(m)
	r, err := client.GetRequestByID(d.Id())
	if err != nil {
		return false, err
	}
	if r.RestyResponse.StatusCode() == 404 {
		//return false, fmt.Errorf("%v", r.RestyResponse.Status())
		return false, nil
	}
	return true, nil

}

func createRequest(d *schema.ResourceData, m interface{}) error {
	id := strings.ToLower(d.Get("request_id").(string))
	numRetriesOnFailure := int64(d.Get("num_retries_on_failure").(int))
	cronFormat := d.Get("schedule").(string)
	scheduleType := strings.ToUpper(d.Get("schedule_type").(string))
	requestType := strings.ToLower(d.Get("request_type").(string))
	instances := int64(d.Get("instances").(int))
	slavePlacement := strings.ToUpper(d.Get("slave_placement").(string))

	// This is a workaround when Singularity delete existing object. Takes several
	// seconds normally.
	timeout := 30
	log.Printf("[TRACE] WAITING for %d seconds", timeout)
	time.Sleep(time.Duration(timeout) * time.Second)

	// Singularity expects uppercase of these values and in our validator,
	// we expect only uppercase to make our resource simpler. Having said
	// that, it does not hurt to always check for value/s in same lowercase.
	log.Printf("Singularity request  '%s' is being provisioned...", id)
	if requestType == "run_once" {
		resp, err := singularity.NewRequest(singularity.RUN_ONCE, id).
			SetInstances(instances).
			SetNumRetriesOnFailures(numRetriesOnFailure).
			SetSlavePlacement(slavePlacement).
			Create(clientConn(m))
		return checkResponse(d, m, resp, err)
	}
	if requestType == "scheduled" {
		req := singularity.NewRequest(singularity.SCHEDULED, "")
		_, err := req.SetScheduleType(scheduleType)
		if err != nil {
			return fmt.Errorf("scheduleType invalid: %v", err)
		}
		_, err = req.SetSchedule(cronFormat)
		if err != nil {
			return fmt.Errorf("cronFormat invalid: %v", err)
		}

		if instances > 1 {
			return fmt.Errorf("Scheduled request can only have instance of: %d", 1)
		}
		resp, err := req.SetNumRetriesOnFailures(numRetriesOnFailure).
			SetID(id).
			SetNumRetriesOnFailures(numRetriesOnFailure).
			SetSlavePlacement(slavePlacement).
			Create(clientConn(m))

		if err != nil {
			return fmt.Errorf("Create new scheduled type request error %v", err)
		}
		return checkResponse(d, m, resp, err)
	}
	if requestType == "service" {
		resp, err := singularity.NewRequest(singularity.SERVICE, id).
			SetInstances(instances).
			SetSlavePlacement(slavePlacement).
			Create(clientConn(m))
		return checkResponse(d, m, resp, err)
	}
	if requestType == "on_demand" {
		resp, err := singularity.NewRequest(singularity.ON_DEMAND, id).
			SetNumRetriesOnFailures(numRetriesOnFailure).
			SetInstances(instances).
			SetSlavePlacement(slavePlacement).
			Create(clientConn(m))
		return checkResponse(d, m, resp, err)
	}
	if requestType == "worker" {
		resp, err := singularity.NewRequest(singularity.WORKER, id).
			SetInstances(instances).
			SetSlavePlacement(slavePlacement).
			Create(clientConn(m))
		return checkResponse(d, m, resp, err)
	}

	return nil
}

func checkResponse(d *schema.ResourceData, m interface{}, r singularity.HTTPResponse, err error) error {
	log.Printf("[TRACE] HTTP Response %v", r.RestyResponse)

	if err != nil {
		return fmt.Errorf("Create Singularity request error: %v", err)
	}
	if r.RestyResponse.StatusCode() <= 200 && r.RestyResponse.StatusCode() >= 299 {
		return fmt.Errorf("Create Singularity request error %v: %v", r.RestyResponse.StatusCode(), err)
	}
	return resourceRequestRead(d, m)
}

// resourceRequestRead is called to resync the local state with the remote state.
// Terraform guarantees that an existing ID will be set. This ID should be used
// to look up the resource. Any remote data should be updated into the local data.
// No changes to the remote resource are to be made.
func resourceRequestRead(d *schema.ResourceData, m interface{}) error {
	client := clientConn(m)
	r, err := client.GetRequestByID(d.Id())
	if err != nil {
		return err
	}
	if r.RestyResponse.StatusCode() == 404 {
		return fmt.Errorf("%v", string(r.RestyResponse.Body()))
	}
	d.Set("request_id", r.Body.SingularityRequest.ID)
	d.Set("request_type", r.Body.SingularityRequest.RequestType)
	d.Set("slave_placement", r.Body.SingularityRequest.SlavePlacement)

	// Only these three types of request expects instance number set.
	if checkRequestTypeMatch(r.Body, "ON_DEMAND", "WORKER", "SERVICE") {
		d.Set("instances", r.Body.SingularityRequest.Instances)
	}
	// Only a scheuled type service expect below parameters.
	if checkRequestTypeMatch(r.Body, "SCHEDULED") {
		d.Set("schedule", r.Body.SingularityRequest.Schedule)
		d.Set("schedule_type", r.Body.SingularityRequest.ScheduleType)
	}

	// Only a service or run_once or on_demand type expect below parameters.
	if checkRequestTypeMatch(r.Body, "SCHEDULED", "RUN_ONCE", "ON_DEMAND") {
		d.Set("num_retries_on_failure", r.Body.SingularityRequest.NumRetriesOnFailure)
	}
	return nil
}

func resourceRequestUpdate(d *schema.ResourceData, m interface{}) error {

	if d.HasChange("request_id") ||
		d.HasChange("schedule") ||
		d.HasChange("request_type") ||
		d.HasChange("num_retries_on_failure") ||
		d.HasChange("schedule_type") ||
		d.HasChange("slave_placement") {
		log.Printf("[TRACE] Delete and update existing request id (%s) success", d.Id())
		// TODO: Investigate whether we can just update existing request, rather
		// than delete and add.
		// Delete existing request and add if there are changes. I couldn't manage
		// to find API doco to update existing request. Only for existing deploy.
		err := resourceRequestDelete(d, m)
		if err != nil {
			return err
		}
		return resourceRequestCreate(d, m)
	}
	if d.HasChange("instances") {
		return resourceScaleRequest(d, m)
	}
	return nil
}

func resourceRequestDelete(d *schema.ResourceData, m interface{}) error {
	a := deleteRequest(d.Id())
	return a(d, m)
}

func deleteRequest(id string) (f func(d *schema.ResourceData, m interface{}) error) {
	return func(d *schema.ResourceData, m interface{}) error {
		req := singularity.NewDeleteRequest(id,
			"Terraform detected changes",
			"Terraform update",
			false)
		resp, err := singularity.DeleteRequest(clientConn(m), req)
		if err != nil {
			return err
		}
		if resp.RestyResponse.StatusCode() == 404 {
			// This could have been deleted manually.
			return nil
		}
		d.SetId("")
		return nil
	}
}

func resourceResourceRequestImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceRequestRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func checkRequestTypeMatch(r singularity.Request, services ...string) bool {
	for _, service := range services {
		if strings.ToUpper(r.RequestType) == strings.ToUpper(service) {
			return true
		}
	}
	return false
}
