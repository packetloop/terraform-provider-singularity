package singularity

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	singularity "github.com/lenfree/go-mesos-singularity"
)

func resourceRequest() *schema.Resource {
	return &schema.Resource{
		Create: resourceRequestCreate,
		Read:   resourceRequestRead,
		Update: resourceRequestUpdate,
		Delete: resourceRequestDelete,

		Schema: map[string]*schema.Schema{
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"request_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"num_retries_on_failure": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instances": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourceRequestCreate(d *schema.ResourceData, m interface{}) error {

	id := d.Get("request_id").(string)
	d.SetId(id)
	return createRequest(d, m)
}

func createRequest(d *schema.ResourceData, m interface{}) error {
	id := d.Get("request_id").(string)
	numRetriesOnFailure := int64(d.Get("num_retries_on_failure").(int))
	cronFormat := d.Get("schedule").(string)
	scheduleType := d.Get("schedule_type").(string)
	requestType := d.Get("request_type").(string)
	instances := int64(d.Get("instances").(int))
	client := m.(*Conn).sclient

	if requestType == "RUN_ONCE" {
		req := singularity.NewRunOnceRequest(id, instances)
		return checkResponse(singularity.CreateRequest(client, req))
	}
	if requestType == "SCHEDULED" {
		req, err := singularity.NewScheduledRequest(id, cronFormat, scheduleType)
		if err != nil {
			return fmt.Errorf("Create new scheduled type request error %v", err)
		}
		req.NumRetriesOnFailure = numRetriesOnFailure
		return checkResponse(singularity.CreateRequest(client, req))
	}
	return nil
}
func checkResponse(r singularity.HTTPResponse, err error) error {
	log.Printf("[TRACE] HTTP Response %v", r.GoRes)

	if err != nil {
		return fmt.Errorf("Create Singularity request error: %v", err)
	}
	if r.GoRes.StatusCode <= 200 && r.GoRes.StatusCode >= 299 {
		return fmt.Errorf("Create Singularity request error %v: %v", r.GoRes.StatusCode, err)
	}

	return nil
}

func resourceRequestRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceRequestUpdate(d *schema.ResourceData, m interface{}) error {
	// Enable partial state mode
	d.Partial(true)

	if d.HasChange("id") {
		// Try updating the address
		if err := updateAddress(d, m); err != nil {
			return err
		}

		d.SetPartial("id")
	}

	// If we were to return here, before disabling partial mode below,
	// then only the "address" field would be saved.

	// We succeeded, disable partial mode. This causes Terraform to save
	// save all fields again.
	d.Partial(false)

	return nil
}

func updateAddress(d *schema.ResourceData, m interface{}) error {
	return nil //fmt.Errorf("Error update address")
}

func resourceRequestDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
