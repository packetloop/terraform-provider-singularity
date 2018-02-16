package singularity

import (
	"testing"
)

func TestOnDemandRequestDefault(t *testing.T) {
	var data = []struct {
		expectedID   string
		expectedType string
	}{
		{"test-id", "ON_DEMAND"},
		{"demand-123", "ON_DEMAND"},
	}

	for _, tt := range data {
		req := NewRequest(ON_DEMAND, tt.expectedID).Get()
		if req.ID != tt.expectedID {
			t.Errorf("NewRequest(%s, %s): expected %s, got %s", "ON_DEMAND", tt.expectedID, tt.expectedID, req.ID)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("NewRequest(%s, %s): expected %s, got %s", "ON_DEMAND", tt.expectedType, tt.expectedType, req.RequestType)
		}
	}
}

func TestOnDemandRequestSet(t *testing.T) {
	var data = []struct {
		initID       string
		expectedID   string
		expectedType string
	}{
		{"myid", "test-id", "ON_DEMAND"},
		{"odl-idname", "demand-123", "ON_DEMAND"},
		{"", "new-id", "ON_DEMAND"},
	}

	for _, tt := range data {
		req := NewRequest(ON_DEMAND, tt.initID).SetID(tt.expectedID).Get()
		if req.ID != tt.expectedID {
			t.Errorf("NewRequest(%s, %s): expected %s, got %s", "ON_DEMAND", tt.initID, tt.expectedID, req.ID)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("NewRequest(%s, %s): expected %s, got %s", "ON_DEMAND", tt.initID, tt.expectedID, req.RequestType)
		}
	}
}

func TestNewServiceRequestDefault(t *testing.T) {
	var data = []struct {
		expectedID        string
		expectedType      string
		expectedInstances int64
	}{
		{"test-id", "SERVICE", 1},
		{"service-123", "SERVICE", 1},
	}

	for _, tt := range data {
		req := NewRequest(SERVICE, tt.expectedID).Get()
		if req.ID != tt.expectedID {
			t.Errorf("Got %s, expected %s", req.ID, tt.expectedID)
		}
		if req.Instances != tt.expectedInstances {
			t.Errorf("Got %v, expected %v", req.Instances, tt.expectedInstances)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("Got %s, expected %s", req.RequestType, tt.expectedType)
		}
	}
}
func TestNewServiceRequestSet(t *testing.T) {
	var data = []struct {
		expectedID         string
		expectedType       string
		expectedInstances  int64
		expectedNumRetries int64
	}{
		{"test-id", "SERVICE", 0, 5},
		{"test123", "SERVICE", 2, 3},
	}

	for _, tt := range data {
		req := NewRequest(SERVICE, tt.expectedID).
			SetInstances(tt.expectedInstances).
			SetNumRetriesOnFailures(tt.expectedNumRetries).
			Get()
		if req.ID != tt.expectedID {
			t.Errorf("Got %s, expected %s", req.ID, tt.expectedID)
		}
		if req.Instances != tt.expectedInstances {
			t.Errorf("Got %v, expected %v", req.Instances, tt.expectedInstances)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("Got %s, expected %s", req.RequestType, tt.expectedType)
		}
	}
}

func TestNewScheduledRequestSet(t *testing.T) {
	var data = []struct {
		actualID             string
		actualType           string
		actualCron           string
		actualScheduleType   string
		expectedID           string
		expectedType         string
		expectedCron         string
		expectedScheduleType string
		expectedError        bool
	}{
		{"test-scheduled", "SCHEDULED", "*/30 * * * *", "CRON", "test-scheduled", "SCHEDULED", "*/30 * * * *", "CRON", false},
		{"failed-scheduled", "SCHEDULED", "* * * * * * *", "CRON", "failed-scheduled", "SCHEDULED", "Parse * * * * * * * cron schedule error. Expected 5 or 6 fields, found 7: * * * * * * *", "CRON", true},
	}

	for _, tt := range data {
		s := NewRequest(SCHEDULED, tt.actualID)
		sched, _ := s.SetScheduleType(tt.actualScheduleType)
		schedType, err := sched.SetSchedule(tt.actualCron)
		// Catch invalid cron format which should return error.
		if tt.expectedError == true {
			if err == nil {
				t.Errorf("SetSchedule(%s): expected %v, actual %v", tt.actualCron, tt.expectedCron, err.Error())
			}
		}
		if err == nil {
			req := schedType.SetID(tt.actualID).Get()

			if req.ID != tt.expectedID {
				t.Errorf("SetID(%s): expected %v, actual %v", tt.actualID, tt.expectedID, req.ID)
			}
			if req.Schedule != tt.expectedCron {
				t.Errorf("SetSchedule(%s): expected %v, actual %v", tt.actualCron, tt.expectedCron, req.Schedule)
			}
			if req.RequestType != tt.expectedType {
				t.Errorf("NewRequest(%s, %s): expected %v, actual %v", "SCHEDULED", tt.actualID, tt.expectedCron, req.Schedule)
			}
			if req.ScheduleType != tt.expectedScheduleType {
				t.Errorf("SetScheduleType(%s): expected %v, actual %v", tt.actualScheduleType, tt.expectedCron, req.Schedule)
			}
		}
	}
}

func TestNewWorkerRequestDefault(t *testing.T) {
	var data = []struct {
		expectedID        string
		expectedType      string
		expectedInstances int64
	}{
		{"test-id", "WORKER", 1},
		{"test-id-3", "WORKER", 1},
	}

	for _, tt := range data {
		req := NewRequest(WORKER, tt.expectedID).Get()
		if req.ID != tt.expectedID {
			t.Errorf("Got %s, expected %s", req.ID, tt.expectedID)
		}
		if req.Instances != tt.expectedInstances {
			t.Errorf("Got %v, expected %v", req.Instances, tt.expectedInstances)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("Got %s, expected %s", req.RequestType, tt.expectedType)
		}
	}
}
func TestNewWorkerRequestSet(t *testing.T) {
	var data = []struct {
		expectedID        string
		expectedType      string
		expectedInstances int64
	}{
		{"test-id", "WORKER", 0},
		{"test-id-2", "WORKER", 2},
	}

	for _, tt := range data {
		req := NewRequest(WORKER, tt.expectedID).SetInstances(tt.expectedInstances).Get()
		if req.ID != tt.expectedID {
			t.Errorf("Got %s, expected %s", req.ID, tt.expectedID)
		}
		if req.Instances != tt.expectedInstances {
			t.Errorf("Got %v, expected %v", req.Instances, tt.expectedInstances)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("Got %s, expected %s", req.RequestType, tt.expectedType)
		}
	}
}

func TestNewRunOnceRequestDefault(t *testing.T) {
	var data = []struct {
		expectedID        string
		expectedType      string
		expectedInstances int64
	}{
		{"test-id", "RUN_ONCE", 1},
		{"test-id-2", "RUN_ONCE", 1},
	}

	for _, tt := range data {
		req := NewRequest(RUN_ONCE, tt.expectedID).Get()
		t.Logf("%v", req)
		if req.ID != tt.expectedID {
			t.Errorf("Got %s, expected %s", req.ID, tt.expectedID)
		}
		if req.Instances != tt.expectedInstances {
			t.Errorf("Got %v, expected %v", req.Instances, tt.expectedInstances)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("Got %s, expected %s", req.RequestType, tt.expectedType)
		}
	}
}

func TestNewRunOnceSet(t *testing.T) {
	var data = []struct {
		expectedID        string
		expectedType      string
		expectedInstances int64
	}{
		{"test-id", "RUN_ONCE", 0},
		{"test-id-2", "RUN_ONCE", 2},
	}

	for _, tt := range data {
		req := NewRequest(RUN_ONCE, tt.expectedID).SetInstances(tt.expectedInstances).Get()
		if req.ID != tt.expectedID {
			t.Errorf("Got %s, expected %s", req.ID, tt.expectedID)
		}
		if req.Instances != tt.expectedInstances {
			t.Errorf("Got %v, expected %v", req.Instances, tt.expectedInstances)
		}
		if req.RequestType != tt.expectedType {
			t.Errorf("Got %s, expected %s", req.RequestType, tt.expectedType)
		}
	}
}

func TestNewRequestScale(t *testing.T) {
	expectedID := "scale-id-test"
	expectedInstances := 3
	expectedMessage := "test scale"
	expectedIncrement := 2
	req := NewRequestScale(expectedID,
		expectedMessage,
		expectedInstances,
		expectedIncrement)
	if req.id != expectedID {
		t.Errorf("Got %s, expected %s", req.id, expectedID)
	}
	if req.SingularityScaleRequest.Instances != expectedInstances {
		t.Errorf("Got %v, expected %v ",
			req.SingularityScaleRequest.Instances,
			expectedInstances)
	}
	if req.SingularityScaleRequest.Message != expectedMessage {
		t.Errorf("Got %v, expected %v ",
			req.SingularityScaleRequest.Message,
			expectedMessage)
	}
	if req.SingularityScaleRequest.Incremental != expectedIncrement {
		t.Errorf("Got %v, expected %v ",
			req.SingularityScaleRequest.Incremental,
			expectedIncrement)
	}
}

func TestNewDeploy(t *testing.T) {
	expectedRequestID := "test-id-1"
	expectedDeployID := "4"

	req := NewDeleteDeploy(expectedRequestID, expectedDeployID)
	if req.deployID != expectedDeployID {
		t.Errorf("Got %v, expected %v ", req.deployID, expectedDeployID)
	}
	if req.requestID != expectedRequestID {
		t.Errorf("Got %v, expected %v ", req.requestID, expectedRequestID)
	}
}

/*  Fix this http request test. Checkout gomega http client test
https://onsi.github.io/gomega/#ghttp-testing-http-clients

func TestClient_GetRequests(t *testing.T) {
       request := SingularityRequest{
               ID:                  "test-geostreamoffsets-launch-sqs-connector",
               requestType:         "RUN_ONCE",
               NumRetriesOnFailure: 3,
       }
       activeDeploy := ActiveDeploy{
               RequestID: "test-geostreamoffsets-launch-sqs-connector",
               DeployID:  "prodromal",
               Timestamp: 1503451301091,
       }
       deployState := SingularityDeployState{
               RequestID:    "test-geostreamoffsets-launch-sqs-connector",
               ActiveDeploy: activeDeploy,
       }
       data := Requests{
               Request{
                       SingularityRequest: request,
                       State:              "ACTIVE",
                       SingularityDeployState: deployState,
               },
       }

       config := Config{
               Host: "127.0.0.1",
       }
       c := New(config)

       httpmock.Activate()
       defer httpmock.DeactivateAndReset()
       da, _ := json.Marshal(data)
       httpmock.NewMockTransport().RegisterResponder("GET", "http://foo.com/bar", httpmock.NewStringResponder(200, string(da)))

       req, _, _ := c.SuperAgent.Get("http://foo.com/bar").End()
       //      req, _, _ := c.GetRequests()
       //      req, _ := http.NewRequest("GET", "http://foo.com/bar", nil)

       fmt.Println("val: ", req)
       //res, _ := (&http.Client{}).Do(req)
       z, _ := ioutil.ReadAll(req.Body)
       fmt.Println("val: ", string(z))

               st.Expect(t, err, nil)
               st.Expect(t, res.StatusCode, 200)

               // Verify that we don't have pending mocks
               st.Expect(t, gock.IsDone(), true)
}
*/
