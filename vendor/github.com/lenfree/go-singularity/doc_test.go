// +build !doc

package singularity_test

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	singularity "github.com/lenfree/go-singularity"
)

func ExampleCreateRequest() {

	const (
		ON_DEMAND = 1
		SERVICE   = 2
		SCHEDULED = 3
		RUN_ONCE  = 4
		WORKER    = 5
	)

	c := singularity.NewConfig().
		SetHost("localhost/singularity").
		SetPort(80).
		SetRetry(3).
		Build()
	client := singularity.NewClient(c)
	res, _ := singularity.NewRequest(ON_DEMAND, "").SetID("lenfree-test").Create(client)
	fmt.Println(res.RestyResponse.Status())
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, _ := json.Marshal(res.Body)
	fmt.Println(string(data))

	// Output:
	// {"request":
	//	{
	//		"id":"lenfree-test",
	//		"requestType":"ON_DEMAND",
	//		"numRetriesOnFailure":0,
	//		"rackSensitive":false,
	//		"loadBalanced":false,
	//		"killOldNonLongRunningTasksAfterMillis":0,
	//		"scheduledExpectedRuntimeMillis":0,
	//		"bounceAfterScale":false,
	//		"skipHealthchecks":false,
	//		"taskLogErrorRegex": "",
	//		"taskLogErrorRegexCaseSensitive":false
	//	},
	//	"state":"ACTIVE"}"
	//)
	// Output:
	// 200 OK
}

func ExampleClient_GetRequestByID() {
	c := singularity.NewConfig().
		SetHost("localhost/singularity").
		SetPort(80).
		SetRetry(3).
		Build()
	client := singularity.NewClient(c)
	_, r, _ := client.GetRequests()

	// This requestID have a deploy attach to it. Hence,
	// it can be decode to type Task.
	resp, _ := client.GetRequestByID(r[0].ID)
	fmt.Printf("debug: %s\n", resp.Body.ActiveDeploy.ContainerInfo.Docker.Image)

	// Output:
	// golang:latest
}

func ExampleDeleteRequest() {
	c := singularity.NewConfig().
		SetHost("localhost/singularity").
		SetPort(80).
		SetRetry(3).
		Build()
	client := singularity.NewClient(c)
	d := singularity.NewDeleteRequest("lenfree-test-run-once", "test delete", "", false)
	r, _ := singularity.DeleteRequest(client, d)
	fmt.Println(r.Response.ID)

	// Output:
	// lenfree-test-run-once
}

func ExampleScaleRequest() {
	c := singularity.NewConfig().
		SetHost("localhost/singularity").
		SetPort(80).
		SetRetry(3).
		Build()
	client := singularity.NewClient(c)
	s := singularity.NewRequestScale("lenfree-test-run-once", "scale up to 2 by 1 increment", 2, 1)
	r, e := singularity.ScaleRequest(client, *s)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	fmt.Printf("%#v\n", r.RequestParent.SingularityRequest)
	fmt.Printf("expiring: %#v\n", r.RequestParent.SingularityExpiringScale)

	//state: singularity.SingularityRequest{
	//	ID:"lenfree-test-run-once",
	//	Instances:3,
	//	NumRetriesOnFailure:0,
	//	QuartzSchedule:"",
	//	RequestType:"RUN_ONCE",
	//	Schedule:"",
	//	ScheduleType:"",
	//	HideEvenNumberAcrossRacksHint:false,
	//	TaskExecutionTimeLimitMillis:0,
	//	TaskLogErrorRegexCaseSensitive:false,
	//	SkipHealthchecks:false,
	//	WaitAtLeastMillisAfterTaskFinishesForReschedule:0,
	//	TaskPriorityLevel:0,
	//	RackAffinity:[]string(nil),
	//	MaxTasksPerOffer:0,
	//	BounceAfterScale:false,
	//	RackSensitive:false,
	//	AllowedSlaveAttributes:map[string]string(nil),
	//	Owners:[]string(nil),
	//	RequiredRole:"",
	//	ScheduledExpectedRuntimeMillis:0,
	//	RequiredSlaveAttributes:map[string]string(nil),
	//	LoadBalanced:false,
	//	KillOldNonLongRunningTasksAfterMillis:0,
	//	ScheduleTimeZone:"",
	//	AllowBounceToSameHost:false,
	//	TaskLogErrorRegex:""
	//}
	//expiring: singularity.SingularityExpiringScale{
	//	RevertToInstances:2,
	//	User:"",
	//	RequestID:"lenfree-test-run-once",
	//	Bounce:false,
	//	StartMillis:1511057602985,
	//	ActionID:"",
	//	DurationMillis:0,
	//	SingularityExpiringAPIRequestObject:singularity.SingularityExpiringAPIRequestObject{
	//		ActionID:"",
	//		DurationMillis:0,
	//		Instances:3,
	//		Message:"",
	//		SkipHealthchecks:false
	//		}
	//	}
}
