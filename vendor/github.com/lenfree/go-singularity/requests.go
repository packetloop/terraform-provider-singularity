package singularity

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty"
	cron "gopkg.in/robfig/cron.v2"
)

const (
	ON_DEMAND = 1
	SERVICE   = 2
	SCHEDULED = 3
	RUN_ONCE  = 4
	WORKER    = 5
)

// NewRequest accepts 5 types of Singularity reqeuests, a string id
// and return default values for each type of request with request id.
func NewRequest(t int, id string) ServiceRequest {
	switch t {
	case ON_DEMAND:
		return &SingularityRequest{
			RequestType: "ON_DEMAND",
			ID:          id,
		}
	case SERVICE:
		return &SingularityRequest{
			RequestType: "SERVICE",
			Instances:   1,
			ID:          id,
		}
	case SCHEDULED:
		return &SingularityRequest{
			RequestType:  "SCHEDULED",
			ScheduleType: "CRON",
			ID:           id,
		}
	case RUN_ONCE:
		return &SingularityRequest{
			RequestType: "RUN_ONCE",
			Instances:   1,
			ID:          id,
		}
	case WORKER:
		return &SingularityRequest{
			RequestType: "WORKER",
			Instances:   1,
			ID:          id,
		}
	}
	return nil
}

// GetRequests retrieve the list of all Singularity requests.
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#endpoint-/api/requests
func (c *Client) GetRequests() (*resty.Response, Requests, error) {
	var body Requests
	res, err := c.Rest.
		R().
		Get(c.Endpoint + "/api/requests")

	err = c.Rest.JSONUnmarshal(res.Body(), body)
	if err != nil {
		return &resty.Response{}, nil, fmt.Errorf("Get Singularity requests not found: %v", err)
	}
	return res, body, nil
}

// GetRequestByID accpets string id and retrieve a specific Singularity Request by ID
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#get-apirequestsrequestrequestid
func (c *Client) GetRequestByID(id string) (HTTPResponse, error) {
	res, err := c.Rest.
		R().
		Get(c.Endpoint + "/api/requests/request" + "/" + id)

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Get Singularity request not found: %v", err)
	}

	return HTTPResponse{
		RestyResponse: res,
	}, nil
}

// HTTPResponse contains response and body from a http query.
// TODO: Move different response type to use interface{} rather
// than user defined types.
type HTTPResponse struct {
	RestyResponse *resty.Response
	Body          Request
	Task          interface{}
	Response      SingularityRequest
	RequestParent SingularityRequestParent
}

// CreateRequest accepts ServiceRequest struct and Creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
//func CreateRequest(c *Client, r ServiceRequest) (HTTPResponse, error) {
//	return r.Create(c)
//}

// ServiceRequest is an interface to different types of Singularity job requestType.
type ServiceRequest interface {
	Create(*Client) (HTTPResponse, error)
	SetID(string) ServiceRequest
	Get() SingularityRequest
	SetInstances(int64) ServiceRequest
	SetSchedule(string) (ServiceRequest, error)
	SetScheduleType(string) (ServiceRequest, error)
	SetNumRetriesOnFailures(int64) ServiceRequest
}

// SetID accepts a string to assign a request ID.
// This returns a SeVrviceRequest struct which have parameters
// required to Create a ON_DEMAND type of Singularity job/task.
func (r *SingularityRequest) SetID(s string) ServiceRequest {
	r.ID = s
	return r
}

// Get returns ID of a Singularity Request.
func (r *SingularityRequest) Get() SingularityRequest {
	return *r
}

// SetNumRetriesOnFailures accepts an int64 and sets Service request type
// number retires on failures.
func (r *SingularityRequest) SetNumRetriesOnFailures(i int64) ServiceRequest {
	r.NumRetriesOnFailure = i
	return r
}

// SetScheduleType accepts a cron schedule format.
// Only cron is accepted because this is widely know than quartz.
func (r *SingularityRequest) SetScheduleType(t string) (ServiceRequest, error) {
	if strings.ToLower(t) != "cron" {
		return nil, fmt.Errorf("%v", "Only cron scheduleType is allowed.")
	}
	r.ScheduleType = t
	return r, nil
}

// SetSchedule accepts a cron schedule format as string
// and set shedule for this request.
func (r *SingularityRequest) SetSchedule(s string) (ServiceRequest, error) {
	// Singularity Request expects CRON schedule a string. Hence, we just use cron package
	// to parse and validate this value.
	_, err := cron.Parse(s)

	if err != nil {
		return nil, fmt.Errorf("Parse %s cron schedule error. %v", s, err)
	}
	r.Schedule = s
	return r, nil
}

// SetInstances accepts a cron schedule format as string
// and set shedule for this request.
func (r *SingularityRequest) SetInstances(i int64) ServiceRequest {
	r.Instances = i
	return r
}

// Create accepts ServiceRequest struct and Creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *SingularityRequest) Create(c *Client) (HTTPResponse, error) {
	res, err := c.Rest.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(r).
		Post(c.Endpoint + "/api/requests")

	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Create Singularity request error: %v", err)
	}

	var data Request
	err = c.Rest.JSONUnmarshal(res.Body(), data)

	return HTTPResponse{
		RestyResponse: res,
		Body:          data,
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
	res, err := c.Rest.
		R().
		Delete(c.Endpoint + "/api/requests/request/" + r.id)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Delete Singularity request error: %v", err)
	}

	var data SingularityRequest

	err = c.Rest.JSONUnmarshal(res.Body(), &data)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("parse singularity request delete error: %v", err)
	}

	response := HTTPResponse{
		RestyResponse: res,
		Response:      data,
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

// Scale accepts ServiceRequest struct and Creates a Singularity
// job based on a requestType. Valid types are: SERVICE, WORKER, SCHEDULED,
// ON_DEMAND, RUN_ONCE.
func (r *ScaleHTTPRequest) scale(c *Client) (HTTPResponse, error) {
	res, err := c.Rest.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(r.SingularityScaleRequest).
		Put(c.Endpoint + "/api/requests/request/" + r.id + "/scale")
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Scale Singularity request error: %v", err)
	}

	if res.StatusCode() == 400 {
		return HTTPResponse{}, fmt.Errorf("Scale Singularity request error: %v", r.SingularityScaleRequest)
	}
	// TODO: Maybe use interface and type assertion? Since response would have different types
	// of responses based on request body sent.
	var data SingularityRequestParent
	err = c.Rest.JSONUnmarshal(res.Body(), &data)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity request error: %v", err)
	}

	response := HTTPResponse{
		RestyResponse: res,
		RequestParent: data,
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

// Create Creates a deploy and attach to a existing request.
func (r *SingularityDeployRequest) Create(c *Client) (HTTPResponse, error) {
	res, err := c.Rest.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(r).
		Post(c.Endpoint + "/api/deploys/")
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Scale Singularity request error: %v", err)
	}

	if res.StatusCode() == 400 {
		// 400	Deploy object is invalid
		return HTTPResponse{}, fmt.Errorf("Create Singularity deploy error: %v", r)
	}

	// TODO: Maybe use interface and type assertion? Since response would have different types
	// of responses based on request body sent.
	var data SingularityRequestParent
	err = c.Rest.JSONUnmarshal(res.Body(), &data)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity request error: %v", err)
	}
	response := HTTPResponse{
		RestyResponse: res,
		RequestParent: data,
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
	res, err := c.Rest.
		R().
		Delete(c.Endpoint + "/api/deploys/deploy/" + r.deployID + "/request/" + r.requestID)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Delete Singularity deploy  error: %v", err)
	}

	var data SingularityRequestParent

	e := c.Rest.JSONUnmarshal(res.Body(), &data)
	if e != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity deploy delete error: %v", e)
	}

	response := HTTPResponse{
		RestyResponse: res,
		RequestParent: data,
	}
	return response, nil
}
