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
				Default:  3,
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
				ForceNew: true,
			},
			"max_tasks_per_offer": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"state": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRequestState,
			},
		},
	}
}

func resourceRequestCreate(d *schema.ResourceData, m interface{}) error {

	id := d.Get("request_id").(string)
	d.SetId(id)
	return createRequest(d, m)
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
	maxTasksPerOffer := d.Get("max_tasks_per_offer").(int)

	// Singularity expects uppercase of these values and in our validator,
	// we expect only uppercase to make our resource simpler. Having said
	// that, it does not hurt to always check for value/s in same lowercase.
	log.Printf("Singularity request  '%s' is being provisioned...", id)
	if requestType == "run_once" {
		resp, err := singularity.NewRequest(singularity.RUN_ONCE, id).
			SetInstances(instances).
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
			SetInstances(instances).
			SetMaxTasksPerOffer(maxTasksPerOffer).
			Create(clientConn(m))

		if err != nil {
			return fmt.Errorf("Create new scheduled type request error %v", err)
		}
		return checkResponse(d, m, resp, err)
	}
	if requestType == "service" {
		resp, err := singularity.NewRequest(singularity.SERVICE, id).
			SetInstances(instances).
			SetMaxTasksPerOffer(maxTasksPerOffer).
			Create(clientConn(m))
		return checkResponse(d, m, resp, err)
	}
	if requestType == "on_demand" {
		resp, err := singularity.NewRequest(singularity.ON_DEMAND, id).
			SetInstances(instances).
			SetMaxTasksPerOffer(maxTasksPerOffer).
			Create(clientConn(m))
		return checkResponse(d, m, resp, err)
	}
	if requestType == "worker" {
		resp, err := singularity.NewRequest(singularity.WORKER, id).
			SetInstances(instances).
			SetMaxTasksPerOffer(maxTasksPerOffer).
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
	return nil
}

func resourceRequestUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	if d.HasChange("request_id") ||
		d.HasChange("schedule") ||
		d.HasChange("request_type") ||
		d.HasChange("num_retries_on_failure") ||
		d.HasChange("schedule") ||
		d.HasChange("instances") ||
		d.HasChange("schedule_type") ||
		d.HasChange("max_tasks_pe_offer") {
		log.Printf("[TRACE] Delete and update existing request id (%s) success", d.Id())
		// TODO: Investigate whether we can just update existing request, rather
		// than delete and add.
		// Delete existing request and add if there are changes. I couldn't manage
		// to find API doco to update existing request. Only for existing deploy.
		err := resourceRequestDelete(d, m)
		if err != nil {
			return err
		}
		// This is a workaround when Singularity delete existing object. Takes a few seconds
		// normally.
		time.Sleep(5 * time.Second)
		d.Partial(false)
		return resourceRequestCreate(d, m)
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
			return fmt.Errorf("Singularity request ID %v not found", id)
		}
		d.SetId("")
		return nil
	}
}
