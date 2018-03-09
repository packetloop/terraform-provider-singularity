package singularity

import (
	"reflect"
	"testing"
)

func TestNewRequestNil(t *testing.T) {
	var data = []struct {
		id             string
		expectedResult string
	}{
		{"test-id", "nil"},
		{"demand-123", "nil"},
	}

	for _, tt := range data {
		req := NewRequest(6, tt.id)
		if req != nil {
			t.Errorf("NewRequest(%s, %s): expected %s, got %s",
				"ON_DEMAND",
				tt.id,
				tt.expectedResult,
				req)
		}
	}
}

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

func TestNewDeployRequest(t *testing.T) {
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

func TestNewDeleteDeploy(t *testing.T) {
	type args struct {
		requestID string
		deployID  string
	}
	tests := []struct {
		name string
		args args
		want DeleteHTTPDeploy
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeleteDeploy(tt.args.requestID, tt.args.deployID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeleteDeploy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeploy(t *testing.T) {
	var data = []struct {
		id         string
		expectedID string
	}{
		{"test-id", "test-id"},
		{"demand-123", "demand-123"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).Get()
		if req.ID != tt.expectedID {
			t.Errorf("NewDeploy(%s): expected %s, got %s", tt.id, tt.expectedID, req.ID)
		}
	}
}

func TestDeploySetRequestID(t *testing.T) {
	var data = []struct {
		id                        string
		requestID                 string
		expectedRequestID         string
		expectedContainerInfoType string
	}{
		{"test-id", "request123", "request123", "DOCKER"},
		{"demand-123", "myrequest", "myrequest", "DOCKER"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetRequestID(tt.requestID).Build()
		if req.RequestID != tt.expectedRequestID {
			t.Errorf("SetRequestID(%s): expected %s, got %s",
				tt.requestID,
				tt.expectedRequestID,
				req.RequestID)
		}
		if req.ContainerInfo.Type != tt.expectedContainerInfoType {
			t.Errorf("SetRequestID(%s): expected %s, got %s",
				tt.requestID,
				tt.expectedContainerInfoType,
				req.ContainerInfo.Type)
		}
	}
}

func TestDeploySetContainerInfo(t *testing.T) {
	var data = []struct {
		id                    string
		containerInfo         ContainerInfo
		expectedContainerInfo ContainerInfo
	}{
		{
			"test-id", ContainerInfo{
				DockerInfo: DockerInfo{
					ForcePullImage: false,
					Image:          "golang:latest",
					SingularityDockerParameter: SingularityDockerParameter{
						Key:   "hello",
						Value: "world",
					},
				},
				Type: "DOCKER",
				Volumes: []SingularityVolume{
					SingularityVolume{
						HostPath:      "/tmp",
						ContainerPath: "/tmp",
						Mode:          "rw",
					},
				},
			}, ContainerInfo{
				DockerInfo: DockerInfo{
					ForcePullImage: false,
					Image:          "golang:latest",
					SingularityDockerParameter: SingularityDockerParameter{
						Key:   "hello",
						Value: "world",
					},
				},
				Type: "DOCKER",
				Volumes: []SingularityVolume{
					SingularityVolume{
						HostPath:      "/tmp",
						ContainerPath: "/tmp",
						Mode:          "rw",
					},
				},
			},
		},
		{
			"my-new_deploy", ContainerInfo{
				DockerInfo: DockerInfo{
					ForcePullImage: true,
					Image:          "golang:latest",
					SingularityDockerParameter: SingularityDockerParameter{
						Key:   "test",
						Value: "true",
					},
				},
				Type: "DOCKER",
			}, ContainerInfo{
				DockerInfo: DockerInfo{
					ForcePullImage: true,
					Image:          "golang:latest",
					SingularityDockerParameter: SingularityDockerParameter{
						Key:   "test",
						Value: "true",
					},
				},
				Type: "DOCKER",
			},
		},
	}

	for _, tt := range data {
		deploy, _ := NewDeploy(tt.id).SetContainerInfo(tt.containerInfo)
		req := deploy.Build()
		if req.ContainerInfo.DockerInfo.Image != tt.expectedContainerInfo.DockerInfo.Image {
			t.Errorf("SetContainer(%v): expected %v, got %v",
				tt.containerInfo,
				tt.expectedContainerInfo,
				req.ContainerInfo)
		}
		if req.ContainerInfo.DockerInfo.ForcePullImage != tt.expectedContainerInfo.DockerInfo.ForcePullImage {
			t.Errorf("SetContainer(%v): expected %v, got %v",
				tt.containerInfo,
				tt.expectedContainerInfo,
				req.ContainerInfo)
		}
		if req.ContainerInfo.DockerInfo.SingularityDockerParameter.Key != tt.expectedContainerInfo.DockerInfo.SingularityDockerParameter.Key {
			t.Errorf("SetContainer(%v): expected %v, got %v",
				tt.containerInfo,
				tt.expectedContainerInfo.DockerInfo.SingularityDockerParameter.Key,
				req.ContainerInfo.DockerInfo.SingularityDockerParameter.Key)
		}
		if req.ContainerInfo.DockerInfo.SingularityDockerParameter.Value != tt.expectedContainerInfo.DockerInfo.SingularityDockerParameter.Value {
			t.Errorf("SetContainer(%v): expected %v, got %v",
				tt.containerInfo,
				tt.expectedContainerInfo.DockerInfo.SingularityDockerParameter.Value,
				req.ContainerInfo.DockerInfo.SingularityDockerParameter.Value)
		}
		if req.ContainerInfo.Type != tt.expectedContainerInfo.Type {
			t.Errorf("SetContainer(%v): expected %v, got %v",
				tt.containerInfo,
				tt.expectedContainerInfo.Type,
				req.ContainerInfo.Type)
		}
		for i := range req.ContainerInfo.Volumes {
			if req.ContainerInfo.Volumes[i].Mode != tt.expectedContainerInfo.Volumes[i].Mode {
				t.Errorf("SetContainer(%v): expected %v, got %v",
					tt.containerInfo,
					tt.expectedContainerInfo.Volumes[0].Mode,
					req.ContainerInfo.Volumes[0].Mode)
			}
			if req.ContainerInfo.Volumes[i].HostPath != tt.expectedContainerInfo.Volumes[i].HostPath {
				t.Errorf("SetContainer(%v): expected %v, got %v",
					tt.containerInfo,
					tt.expectedContainerInfo.Volumes[0].HostPath,
					req.ContainerInfo.Volumes[0].HostPath)
			}
			if req.ContainerInfo.Volumes[i].ContainerPath != tt.expectedContainerInfo.Volumes[i].ContainerPath {
				t.Errorf("SetContainer(%v): expected %v, got %v",
					tt.containerInfo,
					tt.expectedContainerInfo.Volumes[0].ContainerPath,
					req.ContainerInfo.Volumes[i].ContainerPath)
			}
		}
	}
}
func TestDeploySetURIs(t *testing.T) {
	var data = []struct {
		id           string
		uris         []SingularityMesosArtifact
		expectedUris []SingularityMesosArtifact
	}{
		{
			"test-id",
			[]SingularityMesosArtifact{
				SingularityMesosArtifact{
					URI: "file:///etc/docker.tar.gz",
				},
				SingularityMesosArtifact{
					URI: "file:///etc/instance/environment.json",
				},
			},
			[]SingularityMesosArtifact{
				SingularityMesosArtifact{
					URI: "file:///etc/docker.tar.gz",
				},
				SingularityMesosArtifact{
					URI: "file:///etc/instance/environment.json",
				},
			},
		},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetURIs(tt.uris).Build()
		if len(req.Uris) != len(tt.expectedUris) {
			t.Errorf("SetURIs(%v): expected %v, got %v",
				tt.uris,
				len(tt.expectedUris),
				len(req.Uris))
		}
		if len(tt.expectedUris) > 1 {
			for i := range req.Uris {
				if req.Uris[i].URI != tt.expectedUris[i].URI {
					t.Errorf("SetURIs(%v): expected %s, got %s",
						tt.uris,
						tt.expectedUris[i].URI,
						req.Uris[i].URI)
				}
			}
		}
	}
}

func TestDeploySetResources(t *testing.T) {
	var data = []struct {
		id                string
		resources         SingularityDeployResources
		expectedResources SingularityDeployResources
	}{
		{
			"test-id",
			SingularityDeployResources{
				Cpus:     0.5,
				MemoryMb: 128,
				NumPorts: 1,
			},
			SingularityDeployResources{
				Cpus:     0.5,
				MemoryMb: 128,
				NumPorts: 1,
			},
		}, {
			"mydeploy",
			SingularityDeployResources{
				Cpus:     1.5,
				MemoryMb: 512,
				NumPorts: 3,
			},
			SingularityDeployResources{
				Cpus:     1.5,
				MemoryMb: 512,
				NumPorts: 3,
			},
		},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetResources(tt.resources).Build()
		if req.SingularityDeployResources.Cpus != tt.expectedResources.Cpus {
			t.Errorf("SetResources(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.Cpus,
				req.SingularityDeployResources.Cpus)
		}
		if req.SingularityDeployResources.MemoryMb != tt.expectedResources.MemoryMb {
			t.Errorf("SetResources(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.MemoryMb,
				req.SingularityDeployResources.MemoryMb)
		}
		if req.SingularityDeployResources.NumPorts != tt.expectedResources.NumPorts {
			t.Errorf("SetResources(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.NumPorts,
				req.SingularityDeployResources.NumPorts)
		}
	}
}

func TestDeploySetAutoAdvanceDeploySteps(t *testing.T) {
	var data = []struct {
		id                     string
		autoDeployStep         bool
		expectedAutoDeployStep bool
	}{
		{"test-id", true, true},
		{"demand-123", false, false},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetAutoAdvanceDeploySteps(tt.autoDeployStep).Build()
		if req.AutoAdvanceDeploySteps != tt.expectedAutoDeployStep {
			t.Errorf("SetAutoAdvanceDeploySteps(%v): expected %v, got %v",
				tt.autoDeployStep,
				tt.expectedAutoDeployStep,
				req.AutoAdvanceDeploySteps)
		}
	}
}

func TestDeploySetServiceBasePath(t *testing.T) {
	var data = []struct {
		id           string
		path         string
		expectedPath string
	}{
		{"test-id", "/index", "/index"},
		{"demand-123", "/", "/"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetServiceBasePath(tt.path).Build()
		if req.ServiceBasePath != tt.expectedPath {
			t.Errorf("SetServiceBasePath(%v): expected %v, got %v",
				tt.path,
				tt.expectedPath,
				req.ServiceBasePath)
		}
	}
}

func TestDeploySetMetadata(t *testing.T) {
	var data = []struct {
		id               string
		metadata         map[string]string
		expectedMetadata map[string]string
	}{
		{
			"test-id",
			map[string]string{
				"bird":  "blue",
				"snake": "green",
				"cat":   "black",
			},
			map[string]string{
				"bird":  "blue",
				"snake": "green",
				"cat":   "black",
			},
		}, {
			"mydeploy",
			map[string]string{
				"test": "true",
			},
			map[string]string{
				"test": "true",
			},
		},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetMetadata(tt.metadata).Build()
		eq := reflect.DeepEqual(req.Metadata, tt.expectedMetadata)
		if !eq {
			t.Errorf("SetMetadata(%v): expected %v, got %v",
				tt.metadata,
				tt.expectedMetadata,
				req.Metadata)
		}
	}
}

func TestDeploySetLabels(t *testing.T) {
	var data = []struct {
		id             string
		labels         map[string]string
		expectedLabels map[string]string
	}{
		{
			"test-id",
			map[string]string{
				"bird":  "blue",
				"snake": "green",
				"cat":   "black",
			},
			map[string]string{
				"bird":  "blue",
				"snake": "green",
				"cat":   "black",
			},
		}, {
			"mydeploy",
			map[string]string{
				"test": "true",
			},
			map[string]string{
				"test": "true",
			},
		},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetLabels(tt.labels).Build()
		eq := reflect.DeepEqual(req.Labels, tt.expectedLabels)
		if !eq {
			t.Errorf("SetLabels(%v): expected %v, got %v",
				tt.labels,
				tt.expectedLabels,
				req.Labels)
		}
	}
}
func TestDeploySetUser(t *testing.T) {
	var data = []struct {
		id           string
		user         string
		expectedUser string
	}{
		{"mydeploy", "jon doe", "jon doe"},
		{"newdeploy_id", "me", "me"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetUser(tt.user).Build()
		if req.User != tt.expectedUser {
			t.Errorf("SetUser(%v): expected %v, got %v",
				tt.user,
				tt.expectedUser,
				req.User)
		}
	}
}
func TestDeploySetDeployStepWaitTimeMs(t *testing.T) {
	var data = []struct {
		id           string
		time         int
		expectedTime int
	}{
		{"mydeploy", 30, 30},
		{"newdeploy_id", 2, 2},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetDeployStepWaitTimeMs(tt.time).Build()
		if req.DeployStepWaitTimeMs != tt.expectedTime {
			t.Errorf("SetDeployStepWaitTimeMs(%v): expected %v, got %v",
				tt.time,
				tt.expectedTime,
				req.DeployStepWaitTimeMs)
		}
	}
}

func TestDeploySetSkipHealthchecksOnDeploy(t *testing.T) {
	var data = []struct {
		id            string
		value         bool
		expectedValue bool
	}{
		{"mydeploy", true, true},
		{"newdeploy_id", false, false},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetSkipHealthchecksOnDeploy(tt.value).Build()
		if req.SkipHealthchecksOnDeploy != tt.expectedValue {
			t.Errorf("SetSkipHealthchecksOnDeploy(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.SkipHealthchecksOnDeploy)
		}
	}
}
func TestDeploySetCommand(t *testing.T) {
	var data = []struct {
		id            string
		value         string
		expectedValue string
	}{
		{"mydeploy", "echo 'hello'", "echo 'hello'"},
		{"newdeploy_id", "hostname", "hostname"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetCommand(tt.value).Build()
		if req.Command != tt.expectedValue {
			t.Errorf("SetCommand(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.Command)
		}
	}
}
func TestDeploySetDeployInstanceCountPerStep(t *testing.T) {
	var data = []struct {
		id            string
		value         int
		expectedValue int
	}{
		{"mydeploy", 2, 2},
		{"newdeploy_id", 5, 5},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetDeployInstanceCountPerStep(tt.value).Build()
		if req.DeployInstanceCountPerStep != tt.expectedValue {
			t.Errorf("SetDeployInstanceCountPerStep(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.DeployInstanceCountPerStep)
		}
	}
}
func TestDeploySetConsiderHealthyAfterRunningForSeconds(t *testing.T) {
	var data = []struct {
		id            string
		value         int64
		expectedValue int64
	}{
		{"mydeploy", 10001, 10001},
		{"newdeploy_id", 300, 300},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetConsiderHealthyAfterRunningForSeconds(tt.value).Build()
		if req.ConsiderHealthyAfterRunningForSeconds != tt.expectedValue {
			t.Errorf("SetConsiderHealthyAfterRunningForSeconds(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.ConsiderHealthyAfterRunningForSeconds)
		}
	}
}

func TestDeploySetMaxTaskRetries(t *testing.T) {
	var data = []struct {
		id            string
		value         int
		expectedValue int
	}{
		{"mydeploy", 11, 11},
		{"newdeploy_id", 3, 3},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetMaxTaskRetries(tt.value).Build()
		if req.MaxTaskRetries != tt.expectedValue {
			t.Errorf("SetMaxTaskRetries(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.MaxTaskRetries)
		}
	}
}

func TestDeploySetEnv(t *testing.T) {
	var data = []struct {
		id          string
		env         map[string]string
		expectedEnv map[string]string
	}{
		{
			"test-id",
			map[string]string{
				"bird":  "blue",
				"snake": "green",
				"cat":   "black",
			},
			map[string]string{
				"bird":  "blue",
				"snake": "green",
				"cat":   "black",
			},
		}, {
			"mydeploy",
			map[string]string{
				"test": "true",
			},
			map[string]string{
				"test": "true",
			},
		},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetEnv(tt.env).Build()
		eq := reflect.DeepEqual(req.Env, tt.expectedEnv)
		if !eq {
			t.Errorf("SetLabels(%v): expected %v, got %v",
				tt.env,
				tt.expectedEnv,
				req.Env)
		}
	}
}

func TestDeploySetVersion(t *testing.T) {
	var data = []struct {
		id            string
		value         string
		expectedValue string
	}{
		{"mydeploy", "release-01", "release-01"},
		{"newdeploy_id", "v1.3", "v1.3"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetVersion(tt.value).Build()
		if req.Version != tt.expectedValue {
			t.Errorf("SetVersion(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.Version)
		}
	}
}

func TestDeploySetDeployHealthTimeoutSeconds(t *testing.T) {
	var data = []struct {
		id            string
		value         int64
		expectedValue int64
	}{
		{"mydeploy", 10001, 10001},
		{"newdeploy_id", 300, 300},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetDeployHealthTimeoutSeconds(tt.value).Build()
		if req.DeployHealthTimeoutSeconds != tt.expectedValue {
			t.Errorf("SetDeployHealthTimeoutSeconds(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.DeployHealthTimeoutSeconds)
		}
	}
}

func TestDeploySetArgs(t *testing.T) {
	var data = []struct {
		id           string
		args         string
		args1        string
		expectedArgs []string
	}{
		{
			"test-id",
			"bird",
			"snake",
			[]string{
				"bird",
				"snake",
			},
		}, {
			"mydeploy",
			"test",
			"true",
			[]string{
				"test",
				"true",
			},
		},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetArgs(tt.args, tt.args1).Build()
		eq := reflect.DeepEqual(req.Arguments, tt.expectedArgs)
		if !eq {
			t.Errorf("SetArgs(%v, %v): expected %v, got %v",
				tt.args,
				tt.args1,
				tt.expectedArgs,
				req.Arguments)
		}
	}
}
func TestDeploySetID(t *testing.T) {
	var data = []struct {
		id         string
		newID      string
		expectedID string
	}{
		{"test-id", "request123", "request123"},
		{"demand-123", "myrequest", "myrequest"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetID(tt.newID).Build()
		if req.ID != tt.expectedID {
			t.Errorf("SetID(%s): expected %s, got %s",
				tt.newID,
				tt.expectedID,
				req.ID)
		}
	}
}

func TestDeploySetCustomerExecutorID(t *testing.T) {
	var data = []struct {
		id            string
		value         string
		expectedValue string
	}{
		{"test-id", "myexecutor-id", "myexecutor-id"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetCustomExecutorID(tt.value).Build()
		if req.CustomExecutorID != tt.expectedValue {
			t.Errorf("SetID(%s): expected %s, got %s",
				tt.value,
				tt.expectedValue,
				req.CustomExecutorID)
		}
	}
}

func TestDeploySetCustomerExecutorSource(t *testing.T) {
	var data = []struct {
		id            string
		value         string
		expectedValue string
	}{
		{"test-id", "/bin/custom-executor", "/bin/custom-executor"},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetCustomExecutorSource(tt.value).Build()
		if req.CustomExecutorSource != tt.expectedValue {
			t.Errorf("SetID(%s): expected %s, got %s",
				tt.value,
				tt.expectedValue,
				req.CustomExecutorSource)
		}
	}
}

func TestDeployRequestAttachRequest(t *testing.T) {
	var data = []struct {
		value         Request
		expectedValue SingularityRequest
	}{
		{
			Request{
				SingularityRequest: SingularityRequest{
					ID: "myrequest-id",
				},
			},
			SingularityRequest{
				ID: "myrequest-id",
			},
		},
	}

	for _, tt := range data {
		req := NewDeployRequest().AttachRequest(tt.value).Build()
		eq := reflect.DeepEqual(req.SingularityRequest, tt.expectedValue)
		if !eq {
			t.Errorf("AttachRequest(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.SingularityRequest)
		}
	}
}
func TestDeployRequestSetMessage(t *testing.T) {
	var data = []struct {
		value         string
		expectedValue string
	}{
		{"this is a test message.", "this is a test message."},
	}

	for _, tt := range data {
		req := NewDeployRequest().SetMessage(tt.value).Build()
		if req.Message != tt.expectedValue {
			t.Errorf("SetMessage(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.Message)
		}
	}
}

func TestDeployRequestSetUnpauseOnSuccessfulDeploy(t *testing.T) {
	var data = []struct {
		value         bool
		expectedValue bool
	}{
		{true, true},
	}

	for _, tt := range data {
		req := NewDeployRequest().SetUnpauseOnSuccessfulDeploy(tt.value).Build()
		if req.UnpauseOnSuccessfulDeploy != tt.expectedValue {
			t.Errorf("SetUnpauseOnSuccessfulDeploy(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue,
				req.UnpauseOnSuccessfulDeploy)
		}
	}
}
func TestDeployRequestAttachDeploy(t *testing.T) {
	var data = []struct {
		value         SingularityDeploy
		expectedValue SingularityDeploy
	}{
		{
			SingularityDeploy{
				ID: "my-test-deploy",
			},
			SingularityDeploy{
				ID: "my-test-deploy",
			},
		},
	}

	for _, tt := range data {
		req := NewDeployRequest().AttachDeploy(&tt.value).Build()
		if req.SingularityDeploy.ID != tt.expectedValue.ID {
			t.Errorf("AttachDeploy(%v): expected %v, got %v",
				tt.value,
				tt.expectedValue.ID,
				req.SingularityDeploy.ID)
		}
	}
}

func TestDeployRunNowRequest(t *testing.T) {
	var data = []struct {
		id                string
		resources         SingularityRunNowRequest
		expectedResources SingularityRunNowRequest
	}{
		{
			"test-id",
			SingularityRunNowRequest{
				SingularityDeployResources: SingularityDeployResources{
					Cpus:     0.5,
					MemoryMb: 128,
					NumPorts: 1,
				},
				SkipHealthchecks: true,
				CommandLineArgs:  []string{"-c test"},
				Message:          "test deploy",
			},
			SingularityRunNowRequest{
				SingularityDeployResources: SingularityDeployResources{
					Cpus:     0.5,
					MemoryMb: 128,
					NumPorts: 1,
				},
				SkipHealthchecks: true,
				CommandLineArgs:  []string{"-c test"},
				Message:          "test deploy",
			},
		},
	}

	for _, tt := range data {
		req := NewDeploy(tt.id).SetSingularityRunNowRequest(tt.resources).Build()
		if req.SingularityRunNowRequest.SingularityDeployResources.Cpus != tt.expectedResources.SingularityDeployResources.Cpus {
			t.Errorf("SetSingularityRunNowRequest(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.SingularityDeployResources.Cpus,
				req.SingularityRunNowRequest.SingularityDeployResources.Cpus)
		}
		if req.SingularityRunNowRequest.SingularityDeployResources.MemoryMb != tt.expectedResources.SingularityDeployResources.MemoryMb {
			t.Errorf("SetSingularityRunNowRequest(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.SingularityDeployResources.MemoryMb,
				req.SingularityRunNowRequest.SingularityDeployResources.MemoryMb)
		}
		if req.SingularityRunNowRequest.SingularityDeployResources.NumPorts != tt.expectedResources.SingularityDeployResources.NumPorts {
			t.Errorf("SetSingularityRunNowRequest(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.SingularityDeployResources.NumPorts,
				req.SingularityRunNowRequest.SingularityDeployResources.NumPorts)
		}
		if req.SingularityRunNowRequest.SkipHealthchecks != tt.expectedResources.SkipHealthchecks {
			t.Errorf("SetSingularityRunNowRequest(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.SkipHealthchecks,
				req.SingularityRunNowRequest.SkipHealthchecks)
		}
		if req.SingularityRunNowRequest.Message != tt.expectedResources.Message {
			t.Errorf("SetSingularityRunNowRequest(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.Message,
				req.SingularityRunNowRequest.Message)
		}
		eq := reflect.DeepEqual(req.SingularityRunNowRequest.CommandLineArgs, tt.expectedResources.CommandLineArgs)
		if !eq {
			t.Errorf("SetSingularityRunNowRequest(%v): expected %v, got %v",
				tt.resources,
				tt.expectedResources.CommandLineArgs,
				req.SingularityRunNowRequest.CommandLineArgs)
		}
	}
}
