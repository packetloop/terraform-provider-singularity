// +build !doc

package singularity_test

import (
	"fmt"
	"os"

	singularity "github.com/lenfree/go-singularity"
	"github.com/mitchellh/mapstructure"
)

func ExampleCreateRequest() {
	config := singularity.Config{
		Host: "localhost/singularity",
	}
	client := singularity.New(config)
	onDemandTypeReq := singularity.NewOnDemandRequest("lenfree-test")
	res, _ := singularity.CreateRequest(client, onDemandTypeReq)
	fmt.Println(res.GoRes.Status)
	fmt.Println(res.Body)

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
	config := singularity.Config{
		Host: "localhost/singularity",
	}
	client := singularity.New(config)
	_, r, _ := client.GetRequests()

	// This requestID have a deploy attach to it. Hence,
	// it can be decode to type Task.
	resp, _ := client.GetRequestByID(r[0].ID)
	val := resp.Task.(map[string]interface{})
	var result singularity.Task
	err := mapstructure.Decode(val, &result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("debug: %+#v\n", result.ActiveDeploy.ContainerInfo.Docker.Image)

	// Output:
	// golang:latest
}

func ExampleDeleteRequest() {
	config := singularity.Config{
		Host: "localhost/singularity",
	}
	client := singularity.New(config)
	d := singularity.NewDeleteRequest("lenfree-test-run-once", "test delete", "", false)
	r, _ := singularity.DeleteRequest(client, d)
	fmt.Println(r.Response.ID)

	// Output:
	// lenfree-test-run-once
}

func ExampleScaleRequest() {
	config := singularity.Config{
		Host: "singularity.staging.mayhem.arbor.net/singularity",
	}
	client := singularity.New(config)
	s := singularity.NewRequestScale("lenfree-test-run-once", 3)
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
