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

	var data Request
	err = c.Rest.JSONUnmarshal(res.Body(), &data)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity request error: %v", err)
	}
	return HTTPResponse{
		RestyResponse: res,
		Body:          data,
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
	SetMaxTasksPerOffer(int) ServiceRequest
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

// SetMaxTasksPerOffer accepts a cron schedule format as string
// and set shedule for this request.
func (r *SingularityRequest) SetMaxTasksPerOffer(i int) ServiceRequest {
	r.MaxTasksPerOffer = i
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
	err = c.Rest.JSONUnmarshal(res.Body(), &data)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity create request error: %v", err)
	}

	return HTTPResponse{
		RestyResponse: res,
		Body:          data,
	}, nil
}

// DeleteHTTPRequest contain id string and *SingularityDeployRequest required
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
		return HTTPResponse{}, fmt.Errorf("parse singularity request delete error: %v", string(res.Body()))
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

// DeployRequest is an interface to create a Singularity Deploy object.
type DeployRequest interface {
	Create(*Client) (HTTPResponse, error)
	AttachRequest(Request) DeployRequest
	SetUnpauseOnSuccessfulDeploy(bool) DeployRequest
	SetMessage(string) DeployRequest
	AttachDeploy(Deploy) DeployRequest
	Build() *SingularityDeployRequest
}

// NewDeployRequest returns an empty DeployRequest struct which you could use to set parameters.
func NewDeployRequest() DeployRequest {
	return new(SingularityDeployRequest)
}

// AttachRequest accepts a Singularity Request object and use this request data for this deploy
// , and update the request on successful deploy.
func (r *SingularityDeployRequest) AttachRequest(s Request) DeployRequest {
	req := s.Get()
	r.SingularityRequest = &req
	return r
}

// SetMessage accepts a string message to show users about this deploy (metadata).
func (r *SingularityDeployRequest) SetMessage(m string) DeployRequest {
	r.Message = m
	return r
}

// SetUnpauseOnSuccessfulDeploy accepts bool. If deploy is successful, also unpause the request.
func (r *SingularityDeployRequest) SetUnpauseOnSuccessfulDeploy(b bool) DeployRequest {
	r.UnpauseOnSuccessfulDeploy = b
	return r
}

// AttachDeploy accepts a Singularity Deploy object, containing all the required details
// about the Deploy.
func (r *SingularityDeployRequest) AttachDeploy(d Deploy) DeployRequest {
	r.SingularityDeploy = *d.Build()
	return r
}

// Build returns a SingularityDeployRequest object.
func (r *SingularityDeployRequest) Build() *SingularityDeployRequest {
	return r
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

	if res.StatusCode() >= 200 && res.StatusCode() <= 299 {
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
	// 400	Deploy object is invalid
	return HTTPResponse{}, fmt.Errorf("Create Singularity deploy error: %v", string(res.Body()))

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

// Delete accepts a *Client and delete a existing deploy. This cancel a pending deployment
// (best effort - the deploy may still succeed or fail).
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#delete-apideploysdeploydeployidrequestrequestid
func (r DeleteHTTPDeploy) Delete(c *Client) (HTTPResponse, error) {
	res, err := c.Rest.
		R().
		Delete(c.Endpoint + "/api/deploys/deploy/" + r.deployID + "/request/" + r.requestID)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("Delete Singularity deploy  error: %v", err)
	}

	var data SingularityRequestParent

	e := c.Rest.JSONUnmarshal(res.Body(), &data)
	if e != nil {
		return HTTPResponse{}, fmt.Errorf("Parse Singularity deploy delete error: %v", string(res.Body()))
	}

	response := HTTPResponse{
		RestyResponse: res,
		RequestParent: data,
	}
	return response, nil
}

// Deploy is an interface to create a Singularity Deploy object.
type Deploy interface {
	Build() *SingularityDeploy
	Get() SingularityDeploy
	SetRequestID(string) Deploy
	SetContainerInfo(ContainerInfo) (Deploy, error)
	SetArgs(...string) Deploy
	SetURIs([]SingularityMesosArtifact) Deploy
	SetResources(SingularityDeployResources) Deploy
	SetCustomExecutorID(string) Deploy
	SetCustomExecutorSource(string) Deploy
	SetAutoAdvanceDeploySteps(bool) Deploy
	SetServiceBasePath(string) Deploy
	SetMetadata(map[string]string) Deploy
	SetLabels(map[string]string) Deploy
	SetUser(string) Deploy
	SetDeployStepWaitTimeMs(int) Deploy
	SetSkipHealthchecksOnDeploy(bool) Deploy
	SetCommand(string) Deploy
	SetDeployInstanceCountPerStep(int) Deploy
	SetConsiderHealthyAfterRunningForSeconds(int64) Deploy
	SetSingularityRunNowRequest(SingularityRunNowRequest) Deploy
	SetMaxTaskRetries(int) Deploy
	SetEnv(map[string]string) Deploy
	SetVersion(string) Deploy
	SetID(string) Deploy
	SetDeployHealthTimeoutSeconds(int64) Deploy
}

// NewDeploy accept a deploy ID string and returns a Singularity deploy object.
func NewDeploy(id string) Deploy {
	return &SingularityDeploy{
		ID: id,
		ContainerInfo: ContainerInfo{
			Type: "DOCKER", // We only support DOCKER at the moment.
		},
	}
}

// Get returns a Singularity Deploy object.
func (d *SingularityDeploy) Get() SingularityDeploy {
	return *d
}

// SetRequestID accepts a string request ID which is associated with this deploy.
// This is required.
func (d *SingularityDeploy) SetRequestID(id string) Deploy {
	d.RequestID = id
	return d
}

// Should I create a Constructor for containerinfo object?

// SetContainerInfo accepts a string request ID which is associated with this deploy.
// This is optional.  Currently only supports DOCKER.
func (d *SingularityDeploy) SetContainerInfo(c ContainerInfo) (Deploy, error) {
	if c.Type != "DOCKER" {
		return nil, fmt.Errorf("Error setting Container Type, %v", "Only supports DOCKER. Please create an issue if you need other than DOCKER.")
	}
	d.ContainerInfo = c
	return d, nil
}

// SetArgs accepts variadic string of command arguments. This is optional.
func (d *SingularityDeploy) SetArgs(s ...string) Deploy {
	for _, i := range s {
		d.Arguments = append(d.Arguments, i)
	}
	return d
}

// SetURIs accepts a list of SingularityMesosArtifact. This list
// of URIs to download before executing the deploy command. This is optional.
func (d *SingularityDeploy) SetURIs(u []SingularityMesosArtifact) Deploy {
	for _, i := range u {
		d.Uris = append(d.Uris, i)
	}
	return d
}

// SetResources accepts a SingularityDeployResources object for this deploy. This
// is optional.
func (d *SingularityDeploy) SetResources(r SingularityDeployResources) Deploy {
	d.SingularityDeployResources = r
	return d
}

// SetCustomExecutorID accepts an ID string as Custom Mesos executor id. This
// is optional.
func (d *SingularityDeploy) SetCustomExecutorID(id string) Deploy {
	d.CustomExecutorID = id
	return d
}

// SetCustomExecutorSource accepts a string as Custom Mesos executor source. This
// is optional.
func (d *SingularityDeploy) SetCustomExecutorSource(s string) Deploy {
	d.CustomExecutorSource = s
	return d
}

// SetAutoAdvanceDeploySteps accepts a bool which sets deploy to automatically
// advance to the next target instance count after deployStepWaitTimeMs seconds.
func (d *SingularityDeploy) SetAutoAdvanceDeploySteps(b bool) Deploy {
	d.AutoAdvanceDeploySteps = b
	return d
}

// SetServiceBasePath accepts a string. The base path for the API exposed
// by the deploy. Used in conjunction with the Load balancer API. This
// is optional.
func (d *SingularityDeploy) SetServiceBasePath(p string) Deploy {
	d.ServiceBasePath = p
	return d
}

// SetMetadata accepts a map of string of string. Map of metadata key/value pairs
// associated with the deployment. This is optional.
func (d *SingularityDeploy) SetMetadata(m map[string]string) Deploy {
	d.Metadata = m
	return d
}

// SetLabels accepts map of string of string. Labels for all tasks
// associated with this deploy. This is optional.
func (d *SingularityDeploy) SetLabels(m map[string]string) Deploy {
	d.Labels = m
	return d
}

// SetUser accepts a string and set tasks as this user. This is optional.
func (d *SingularityDeploy) SetUser(u string) Deploy {
	d.User = u
	return d
}

// SetDeployStepWaitTimeMs accepts an time int to wait this long between
// deploy steps. This is optional.
func (d *SingularityDeploy) SetDeployStepWaitTimeMs(t int) Deploy {
	d.DeployStepWaitTimeMs = t
	return d
}

// SetSkipHealthchecksOnDeploy accepts a bool which allows skipping of
// health checks when deploying. This is optional.
func (d *SingularityDeploy) SetSkipHealthchecksOnDeploy(b bool) Deploy {
	d.SkipHealthchecksOnDeploy = b
	return d
}

// SetCommand accepts a string command to execute for this deployment. This
// is optional.
func (d *SingularityDeploy) SetCommand(c string) Deploy {
	d.Command = c
	return d
}

// SetDeployInstanceCountPerStep accepts a count int. Deploy this many instances at a time.
// This parameter is optional.
func (d *SingularityDeploy) SetDeployInstanceCountPerStep(c int) Deploy {
	d.DeployInstanceCountPerStep = c
	return d
}

// SetConsiderHealthyAfterRunningForSeconds accepts a t int64. Number of seconds that a
// service must be healthy to consider the deployment to be successful. This is optional.
func (d *SingularityDeploy) SetConsiderHealthyAfterRunningForSeconds(t int64) Deploy {
	d.ConsiderHealthyAfterRunningForSeconds = t
	return d
}

// SetSingularityRunNowRequest accepts a SinguarltiyRunNowRequest object. Settings used
// to run this deploy immediately. This is optional.
func (d *SingularityDeploy) SetSingularityRunNowRequest(r SingularityRunNowRequest) Deploy {
	d.SingularityRunNowRequest = &r
	return d
}

// SetMaxTaskRetries accepts a count int allowed at most this many failed
// tasks to be retried before failing the deploy. This is optional.
func (d *SingularityDeploy) SetMaxTaskRetries(c int) Deploy {
	d.MaxTaskRetries = c
	return d
}

// SetEnv accepts a map of string of string. This map of environment
//  variable definitions. This is optional.
func (d *SingularityDeploy) SetEnv(e map[string]string) Deploy {
	d.Env = e
	return d
}

// SetVersion accepts a string for deploy version. This is optional.
func (d *SingularityDeploy) SetVersion(v string) Deploy {
	d.Version = v
	return d
}

// SetID accepts an id string and set this deploy ID.
func (d *SingularityDeploy) SetID(id string) Deploy {
	d.ID = id
	return d
}

// SetDeployHealthTimeoutSeconds accepts a time in seconds int64. This number of seconds
// that Singularity waits for this service to become healthy
// (for it to download artifacts, start running, and optionally pass healthchecks.)
//  This is optional.
func (d *SingularityDeploy) SetDeployHealthTimeoutSeconds(t int64) Deploy {
	d.DeployHealthTimeoutSeconds = t
	return d
}

// Build builds a SingularityDeploy object.
func (d *SingularityDeploy) Build() *SingularityDeploy {
	return d
}
