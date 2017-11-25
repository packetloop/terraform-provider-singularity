package singularity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	cron "gopkg.in/robfig/cron.v2"
)

// GetRequests retrieve the list of all Singularity requests.
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#endpoint-/api/requests
func (c *Client) GetRequests() (gorequest.Response, Requests, error) {
	var body Requests
	res, _, err := c.SuperAgent.Get(c.Endpoint+"/api/requests").
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		EndStruct(&body)

	if err != nil {
		return nil, nil, fmt.Errorf("Get Singularity requests not found: %v", err)
	}
	return res, body, nil
}

// GetRequestByID accpets string id and retrieve a specific Singularity Request by ID
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#get-apirequestsrequestrequestid
func (c *Client) GetRequestByID(id string) (HTTPResponse, error) {
	res, body, err := c.SuperAgent.Get(c.Endpoint+"/api/requests/request"+"/"+id).
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		End()

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Get Singularity request not found: %v", err)
	}

	var data Task
	e := json.Unmarshal([]byte(body), &data)

	if e != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity request get error: %v", e)
	}

	response := HTTPResponse{
		GoRes: res,
		Task:  data,
	}
	return response, nil
}

// HTTPResponse contains response and body from a http query.
type HTTPResponse struct {
	GoRes         gorequest.Response
	Body          Request
	Task          Task
	Response      SingularityRequest
	RequestParent SingularityRequestParent
}

// CreateRequest accepts ServiceRequest struct and creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func CreateRequest(c *Client, r ServiceRequest) (HTTPResponse, error) {
	return r.create(c)
}

// ServiceRequest is an interface to different types of Singularity job requestType.
type ServiceRequest interface {
	create(*Client) (HTTPResponse, error)
}

// NewOnDemandRequest accepts a string id and int number of instances.
// This returns a RequestWorker struct which have parameters
// required to create a ON_DEMAND type of Singularity job/task.
func NewOnDemandRequest(id string) *RequestOnDemand {
	return &RequestOnDemand{
		ID:          id,
		RequestType: "ON_DEMAND",
	}
}

// Create accepts ServiceRequest struct and creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *RequestOnDemand) create(c *Client) (HTTPResponse, error) {
	var body Request
	res, _, err := c.SuperAgent.Post(c.Endpoint+"/api/requests").
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		Send(r).
		EndStruct(&body)

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Create Singularity request error: %v", err)
	}

	return HTTPResponse{
		GoRes: res,
		Body:  body,
	}, nil
}

// NewServiceRequest accepts a string id and int number of instances.
// This returns a RequestWorker struct which have parameters
// required to create a SERVICE type of Singularity job/task.
func NewServiceRequest(id string, i int64) *RequestService {
	return &RequestService{
		ID:          id,
		RequestType: "SERVICE",
		Instances:   i,
	}
}

// Create accepts ServiceRequest struct and creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *RequestService) create(c *Client) (HTTPResponse, error) {
	var body Request
	res, _, err := c.SuperAgent.Post(c.Endpoint+"/api/requests").
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		Send(r).
		EndStruct(&body)

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Create Singularity request error: %v", err)
	}

	return HTTPResponse{
		GoRes: res,
		Body:  body,
	}, nil
}

// NewScheduledRequest accepts a string id, cron schedule format and scheduleType as string.
// Only cron is accepted because this is widely know than quartz. This returns a RequestWorker
// struct which have parameter required to
// create a SCHEDULED type of Singularity job/task.
func NewScheduledRequest(id, s, t string) (*RequestScheduled, error) {

	if strings.ToLower(t) != "cron" {
		return nil, fmt.Errorf("%v", "Only cron scheduleType is allowed.")
	}
	// Singularity Request expects CRON schedule a string. Hence, we just use cron package
	// to parse and validate this value.
	_, err := cron.Parse(s)
	if err != nil {
		return &RequestScheduled{}, fmt.Errorf("Parse %s cron schedule error %v", s, err)
	}
	return &RequestScheduled{
		ID:          id,
		RequestType: "SCHEDULED",
		Schedule:    s,
	}, nil
}

// SetCronSchedule accepts a cron schedule format as string
// and set shedule for this request.
func (r *RequestScheduled) SetCronSchedule(s string) error {
	// Singularity Request expects CRON schedule a string. Hence, we just use cron package
	// to parse and validate this value.
	_, err := cron.Parse(s)

	if err != nil {
		return fmt.Errorf("Parse %s cron schedule error %v", s, err)
	}
	r.Schedule = s
	return nil
}

// Create accepts ServiceRequest struct and creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *RequestScheduled) create(c *Client) (HTTPResponse, error) {
	var body Request
	res, _, err := c.SuperAgent.Post(c.Endpoint+"/api/requests").
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		Send(r).
		EndStruct(&body)

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Create Singularity request error: %v", err)
	}

	return HTTPResponse{
		GoRes: res,
		Body:  body,
	}, nil
}

// NewWorkerRequest accepts a string id and int number of instances.
// This returns a RequestWorker struct which have parameters
// required to create a WORKER type of Singularity job/task.
func NewWorkerRequest(id string, i int64) *RequestWorker {
	return &RequestWorker{
		ID:          id,
		RequestType: "WORKER",
		Instances:   i,
	}
}

// Create accepts ServiceRequest struct and creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *RequestWorker) create(c *Client) (HTTPResponse, error) {
	var body Request
	res, _, err := c.SuperAgent.Post(c.Endpoint+"/api/requests").
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		Send(r).
		EndStruct(&body)

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Create Singularity request error: %v", err)
	}

	return HTTPResponse{
		GoRes: res,
		Body:  body,
	}, nil
}

// NewRunOnceRequest accepts a string id and int number of instances.
// This returns a RequestRunOnce struct which have parameters
// required to create a RUN_ONCE type of Singularity job/task.
func NewRunOnceRequest(id string, i int64) *RequestRunOnce {
	return &RequestRunOnce{
		ID:          id,
		RequestType: "RUN_ONCE",
		Instances:   i,
	}
}

// GetID is a placeho
// Create accepts ServiceRequest struct and creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *RequestRunOnce) create(c *Client) (HTTPResponse, error) {
	var body Request
	res, _, err := c.SuperAgent.Post(c.Endpoint+"/api/requests").
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		Send(r).
		EndStruct(&body)

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Create Singularity request error: %v", err)
	}

	return HTTPResponse{
		GoRes: res,
		Body:  body,
	}, nil
}

// DeleteHTTPRequest contain id string and SingularityDeployRequest required
// parameter to delete a existing request.
type DeleteHTTPRequest struct {
	id string
	SingularityDeleteRequest
}

// ServiceDeleteRequest is an interface that accepts id string
// and *Client.
type ServiceDeleteRequest interface {
	deleteRequestByID(*Client)
}

// NewDeleteRequest accepts a request id string, a bool to deletefromloadbalancer,
// string message and action id to associate with for metadata purposes.
func NewDeleteRequest(id, m, a string, b bool) DeleteHTTPRequest {
	return DeleteHTTPRequest{
		id: id,
		SingularityDeleteRequest: SingularityDeleteRequest{
			DeleteFromLoadBalancer: b,
			Message:                m,
			ActionID:               a,
		},
	}
}

// DeleteRequest accepts id as a string and a type DeleteRequest that
// contains metadata when deleting this Request.
func DeleteRequest(c *Client, r DeleteHTTPRequest) (HTTPResponse, error) {
	return r.delete(c)
}

// DeleteRequest accepts id as a string and a type DeleteRequest that
// contains metadata when deleting this Request. This also deletes any
// deploy attach to this requestID.
func (r DeleteHTTPRequest) delete(c *Client) (HTTPResponse, error) {
	res, body, err := c.SuperAgent.Delete(c.Endpoint+"/api/requests/request/"+r.id).
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		Send(r).
		End()

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Delete Singularity request error: %v", err)
	}

	var data SingularityRequest

	e := json.Unmarshal([]byte(body), &data)
	if e != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity request delete error: %v", e)
	}

	response := HTTPResponse{
		GoRes:    res,
		Response: data,
	}
	return response, nil
}

// ScaleHTTPRequest contains a request id and a body parameter required to scale
// in/out of an Singularity request.
type ScaleHTTPRequest struct {
	id string
	SingularityScaleRequest
}

// ServiceScaleRequest is an interface that accepts a *Client and returns
// a HTTPResponse type and error.
type ServiceScaleRequest interface {
	scale(*Client) (HTTPResponse, error)
}

// ScaleRequest accepts a *Client and ScaleHTTPRequest type to scale
// in/out of an existing Singularity request/task.
func ScaleRequest(c *Client, r ScaleHTTPRequest) (HTTPResponse, error) {
	return r.scale(c)
}

// NewRequestScale accepts an id string and a int and returns a pointer to
// type ScaleRequest which have a minimum required paramters to scale a
// Singularity request.
func NewRequestScale(id, m string, i, in int) *ScaleHTTPRequest {
	return &ScaleHTTPRequest{
		id: id,
		SingularityScaleRequest: SingularityScaleRequest{
			Instances:   i,
			Message:     m,
			Incremental: in,
		},
	}
}

// Scale accepts ServiceRequest struct and creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *ScaleHTTPRequest) scale(c *Client) (HTTPResponse, error) {
	res, data, err := c.SuperAgent.Put(c.Endpoint+"/api/requests/request/"+r.id+"/scale").
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		Send(r.SingularityScaleRequest).
		End()

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Scale Singularity request error: %v", err)
	}
	if res.StatusCode == 400 {
		return HTTPResponse{}, fmt.Errorf("Scale Singularity request error: %v", string(data))
	}
	// TODO: Maybe use interface and type assertion? Since response would have different types
	// of responses based on request body sent.
	var body SingularityRequestParent
	e := json.Unmarshal([]byte(data), &body)
	if e != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity request error: %v", e)
	}
	response := HTTPResponse{
		GoRes:         res,
		RequestParent: body,
	}
	return response, nil
}

func NewDeploy(b bool, u SingularityRequest, d SingularityDeploy, m string) *SingularityDeployRequest {
	return &SingularityDeployRequest{
		UnpauseOnSuccessfulDeploy: b,
		SingularityDeploy:         d,
		SingularityRequest:        u,
		Message:                   m,
	}
}

// Create creates a deploy and attach to a existing request.
func (r *SingularityDeployRequest) create(c *Client) (HTTPResponse, error) {
	res, data, err := c.SuperAgent.Post(c.Endpoint+"/api/deploys/").
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusConflict).
		Send(r).
		End()

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Scale Singularity request error: %v", err)
	}
	if res.StatusCode == 400 {
		// 400	Deploy object is invalid
		return HTTPResponse{}, fmt.Errorf("Create Singularity deploy error: %v", string(data))
	}

	// TODO: Maybe use interface and type assertion? Since response would have different types
	// of responses based on request body sent.
	var body SingularityRequestParent
	e := json.Unmarshal([]byte(data), &body)
	if e != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity request error: %v", e)
	}
	response := HTTPResponse{
		GoRes:         res,
		RequestParent: body,
	}
	return response, nil
}

// NewDeleteDeploy accepts a requestID and deployID string and reutnrs
// DeleteHTTPDeploy struct which have a method delete that cancels
// a pending deploy matching both requestID and deployID.
func NewDeleteDeploy(requestID, deployID string) DeleteHTTPDeploy {
	return DeleteHTTPDeploy{
		requestID: requestID,
		deployID:  deployID,
	}
}

// DeleteHTTPDeploy have a struct of requestID and deployID to be use
// to cancel a pending deploy.
type DeleteHTTPDeploy struct {
	requestID string `json:"requestId"`
	deployID  string `json:"deployId"`
}

// DeleteDeploy accepts a *Client and delete a existing deploy. This cancel a pending deployment
// (best effort - the deploy may still succeed or fail).
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#delete-apideploysdeploydeployidrequestrequestid
func (r DeleteHTTPDeploy) delete(c *Client) (HTTPResponse, error) {
	res, body, err := c.SuperAgent.Delete(c.Endpoint+"/api/deploys/deploy/"+r.deployID+"/request/"+r.requestID).
		Retry(3, 5*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusConflict).
		End()

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Delete Singularity deploy  error: %v", err)
	}

	var data SingularityRequestParent

	e := json.Unmarshal([]byte(body), &data)
	if e != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity deploy delete error: %v", e)
	}

	response := HTTPResponse{
		GoRes:         res,
		RequestParent: data,
	}
	return response, nil
}
