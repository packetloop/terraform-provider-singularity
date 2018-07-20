package singularity

// SingularityRequest contains a high level information of
//  for a single project or deployable item.
type SingularityRequest struct {
	ID                                              string            `json:"id"`
	Instances                                       int64             `json:"instances,omitempty"`
	NumRetriesOnFailure                             int64             `json:"numRetriesOnFailure,omitempty"`
	QuartzSchedule                                  string            `json:"quartzSchedule,omitempty"`
	RequestType                                     string            `json:"requestType"`
	Schedule                                        string            `json:"schedule,omitempty"`
	ScheduleType                                    string            `json:"scheduleType,omitempty"`
	HideEvenNumberAcrossRacksHint                   bool              `json:"hideEventNumerAcrossRacksHint,omitempty"`
	TaskExecutionTimeLimitMillis                    int               `json:"taskExecutionTimeLimitMills"`
	TaskLogErrorRegexCaseSensitive                  bool              `json:"taskLogErrorRegexCaseSensitive"`
	SkipHealthchecks                                bool              `json:"skipHealthchecks"`
	WaitAtLeastMillisAfterTaskFinishesForReschedule int               `json:"waitAtleastMillisAfterTaskFinishesForReschedule"`
	TaskPriorityLevel                               int               `json:"taksPriorityLevel"`
	RackAffinity                                    []string          `json:"RackAffinity"`
	MaxTasksPerOffer                                int               `json:"maxTasksPerOffer,omitempty"`
	BounceAfterScale                                bool              `json:"bounceAfterScale"`
	RackSensitive                                   bool              `json:"rackSensitive"`
	AllowedSlaveAttributes                          map[string]string `json:"allowedSlaveAttributes"`
	Owners                                          []string          `json:"owners"`
	RequiredRole                                    string            `json:"requiredRole,omitempty"`
	ScheduledExpectedRuntimeMillis                  int               `json:"scheduledExpectedRuntimeMillis"`
	RequiredSlaveAttributes                         map[string]string `json:"requiredSlaveAttributes"`
	LoadBalanced                                    bool              `json:"loadBalanced,omitempty"`
	KillOldNonLongRunningTasksAfterMillis           int               `json:"killOldNonLongRunningTasksAfterMillis,omitempty"`
	ScheduleTimeZone                                string            `json:"scheduledTimeZone"`
	AllowBounceToSameHost                           bool              `json:"allowBounceToSamehost"`
	TaskLogErrorRegex                               string            `json:"taskLogErrorRegex"`
	SlavePlacement                                  *string           `json:"slavePlacement"`
}

// ActiveDeploy have a string deployId, requestId and a timestamp.
type ActiveDeploy struct {
	DeployID  string `json:"deployId"`
	RequestID string `json:"requestId"`
	Timestamp int64  `json:"timestamp"`
}

// RequestDeployState contains specific configuration or version
// of the running code for that deployable item
type RequestDeployState struct {
	ActiveDeploy `json:"activeDeploy"`
	RequestID    string `json:"requestId"`
}

// Request struct contains all singularity requests.
// This have a JSON response of /api/requests/request/ID.
type Request struct {
	SingularityRequest `json:"request"`
	RequestDeployState struct {
		ActiveDeploy struct {
			DeployID  string `json:"deployId"`
			RequestID string `json:"requestId"`
			Timestamp int64  `json:"timestamp"`
		} `json:"activeDeploy"`
		PendingDeployState struct {
			DeployID  string `json:"deployId"`
			RequestID string `json:"requestId"`
			Timestamp int64  `json:"timestamp"`
		} `json:"pendingDeploy"`
		RequestID string `json:"requestId"`
	} `json:"requestDeployState"`
	State        string `json:"state"`
	ActiveDeploy struct {
		Arguments                  []string `json:"arguments,omitempy"`
		Command                    string   `json:"command,omitempty"`
		ContainerInfo              `json:"containerInfo"`
		Env                        map[string]string `json:"env"`
		ID                         string            `json:"id"`
		RequestID                  string            `json:"requestId"`
		SingularityDeployResources `json:"resources"`
		Uris                       []SingularityMesosArtifact `json:"uris"`
		Volumes                    []SingularityVolume        `json:"volumes"`
		Metadata                   map[string]string          `json:"metadata"`
	} `json:"activeDeploy"`
	PendingDeploy struct {
		CustomExecutorID           string              `json:"customExecutorId"`
		Volumes                    []SingularityVolume `json:"volumes"`
		SingularityDeployResources `json:"resources"`
		Uris                       []SingularityMesosArtifact `json:"uris"`
		ContainerInfo              `json:"containerInfo"`
		Arguments                  []string          `json:"arguments"`
		TaskEnv                    interface{}       `json:"taskEnv"` //Map[int,Map[string,string]]	Map of environment variable overrides for specific task instances.
		AutoAdvanceDeploySteps     bool              `json:"autoAdvanceDeploySteps"`
		ID                         string            `json:"id"`
		Command                    string            `json:"command"`
		Metadata                   map[string]string `json:"metadata"`
	} `json:"pendingDeploy"`
	RunImmediately struct {
		Resources struct {
			Cpus     int64 `json:"cpus"`
			DiskMb   int64 `json:"diskMb"`
			MemoryMb int64 `json:"memoryMb"`
			NumPorts int64 `json:"numPorts"`
		} `json:"resources"`
		RunAt int64  `json:"runAt"`
		RunID string `json:"runId"`
	} `json:"runImmediately"`
	SkipHealthchecksOnDeploy bool `json:"skipHealthchecksOnDeploy"`
}

// Requests is a slice of Request.
type Requests []Request

type Docker struct {
	ForcePullImage bool              `json:"forcePullImage,omitempty"`
	Image          string            `json:"image,omitempty"`
	Parameters     map[string]string `json:"parameters,omitempty"`
	Privileged     bool              `json:"privileged,omitempty"`
}

// ContainerInfo contains information about a Docker type Singularity
// container type.
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#-singularitycontainerinfo
type ContainerInfo struct {
	DockerInfo `json:"docker"`
	Type       string              `json:"type"` // Allowable values: MESOS, DOCKER. Default is MESOS.
	Volumes    []SingularityVolume `json:"volumes,omitempty"`
}

//https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#model-SingularityDockerInfo
type DockerInfo struct {
	Parameters                  map[string]string            `json:"parameters,omitempty"`
	ForcePullImage              bool                         `json:"forcePullImage,omitempty"`
	SingularityDockerParameters []SingularityDockerParameter `json:"dockerParameters,omitEmpty"`
	Privileged                  bool                         `json:"privileged,omitEmpty"`
	Network                     string                       `json:"network,omitEmpty"` //Value can be BRIDGE, HOST, or NONE
	Image                       string                       `json:"image"`
	PortMappings                []DockerPortMapping          `json:"portMappings,omitempty"`
}

//https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#model-SingularityDockerPortMapping
type DockerPortMapping struct {
	ContainerPort     int64  `json:"containerPort"`
	ContainerPortType string `json:"containerPortType,omitempty"` //Allowable values: LITERAL, FROM_OFFER
	HostPort          int64  `json:"hostPort"`
	HostPortType      string `json:"hostPortType,omitempty"` //Allowable values: LITERAL, FROM_OFFER
	Protocol          string `json:"protocol,omitempty"`     //Default is tcp
}

// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#model-SingularityDockerParameter
type SingularityDockerParameter struct {
	Key   string `json:"key,omitEmpty"`
	Value string `json:"value,omitEmpty"`
}

// SingularityVolume contains information about Docker volume. This is optional.
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#model-SingularityVolume
type SingularityVolume struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
	Mode          string `json:"mode"`
}

// Task contains JSON response of /api/requests/request/ID.
type Task struct {
	ActiveDeploy struct {
		Arguments                  []string `json:"arguments"`
		Command                    string   `json:"command"`
		ContainerInfo              `json:"containerInfo"`
		Env                        map[string]string `json:"env"`
		ID                         string            `json:"id"`
		RequestID                  string            `json:"requestId"`
		SingularityDeployResources `json:"resources"`
		Uris                       []string `json:"uris"`
	} `json:"activeDeploy"`
	RequestDeployState struct {
		ActiveDeploy struct {
			DeployID  string `json:"deployId"`
			RequestID string `json:"requestId"`
			Timestamp int64  `json:"timestamp"`
		} `json:"activeDeploy"`
		RequestID string `json:"requestId"`
	} `json:"requestDeployState"`
	State              string `json:"state"`
	SingularityRequest `json:"request"`
}

// SingularityDeployResources includes information about required/configured
// resources needed for a request/job.
type SingularityDeployResources struct {
	Cpus     float64 `json:"cpus,omitempty"`
	MemoryMb float64 `json:"memoryMb,omitempty"`
	DiskMb   float64 `json:"diskMb,omitempty"`
	NumPorts int64   `json:"numPorts,omitempty"`
}

// SingularityScaleRequest contains parameters for making scaling a request. For more info, please see:
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#-singularityscalerequest
type SingularityScaleRequest struct {
	SkipHealthchecks bool   `json:"skipHealthchecks"`
	DurationMillis   int64  `json:"durationMillis"`
	Bounce           bool   `json:"bounce"`
	Message          string `json:"message"`
	ActionID         string `json:"actionId"`
	Instances        int    `json:"instances"`
	Incremental      int    `json:"incremental"`
}

// SingularityExpiringSkipHealthchecks have parameters for a expiring skip
// healthchecks.
type SingularityExpiringSkipHealthchecks struct {
	User                                string      `json:"user"`
	RequestID                           string      `json:"requestId"`
	StartMillis                         int64       `json:"startMillis"`
	ActionID                            string      `json:"actionId"`
	SingularityExpiringAPIRequestObject interface{} `json:"expiringAPIRequestObject"`
	RevertToSkipHealthchecks            bool        `json:"revertToSkipHealthchecks"`
}

// RequestState contains a string state of a existing job/Singulariy request. Allowable
// values are:  ACTIVE, DELETING, DELETED, PAUSED, SYSTEM_COOLDOWN, FINISHED,
// DEPLOYING_TO_UNPAUSE
type RequestState struct {
	State string `json:"state"` //Allowable values:
}

// HealthcheckProtocol contains a string with allowable value of
// HTTP or HTTPS.
type HealthcheckProtocol string

// HealthcheckOptions contains parameters of a healthcheck options
// for a new and existing Singularity request.
type HealthcheckOptions struct {
	StartupDelaySeconds    int    `json:"startupDelaySeconds"`
	ResponseTimeoutSeconds int    `json:"responseTimeoutSeconds"`
	IntervalSeconds        int    `json:"intervalSeconds"`
	URI                    string `json:"uri"` //Healthcheck uri to hit
	FailureStatusCodes     []int  `json:"failureStatusCodes"`
	MaxRetries             int    `json:"maxRetries"`
	StartupTimeoutSeconds  int    `json:"startupTimeoutSeconds"`
	PortNumber             int    `json:"portNumber"`
	StartupIntervalSeconds int    `json:"startupIntervalSeconds"` //Time to wait after a failed healthcheck to try again in seconds.
	HealthcheckProtocol    `json:"protocol"`
	PortIndex              int `json:"portIndex"`
}

type SingularityMesosTaskLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EmbeddedArtifact struct {
	TargetFolderRelativeToTask string `json:"targetFolderRelativeToTask"` // optional
	Md5sum                     string `json:"md5sum"`                     // optional
	Filename                   string `json:"filename"`                   // optional
	Name                       string `json:"name"`                       // optional
	Content                    []byte `json:"content"`                    // optional
}

type S3Artifact struct {
	TargetFolderRelativeToTask string `json:"targetFolderRelativeToTask"` //	optional
	S3Bucket                   string `json:"s3Bucket"`                   //optional
	Md5sum                     string `json:"md5sum"`                     // optional
	Filename                   string `json:"filename"`                   //optional
	Filesize                   int64  `json:"filesize"`                   // long	// optional
	S3ObjectKey                string `json:"s3ObjectKey"`                // optional
	Name                       string `json:"name"`                       // optional
	IsArtifactList             bool   `json:"isArtifactList"`             //optional
}

type S3ArtifactSignature struct {
	TargetFolderRelativeToTask string `json:"targetFolderRelativeToTask"` //	optional
	S3Bucket                   string `json:"s3Bucket"`                   //optional
	Md5sum                     string `json:"md5sum"`                     // optional
	Filename                   string `json:"filename"`                   //optional
	Filesize                   int64  `json:"filesize"`                   // long	// optional
	S3ObjectKey                string `json:"s3ObjectKey"`                // optional
	Name                       string `json:"name"`                       // optional
	IsArtifactList             bool   `json:"isArtifactList"`             //optional
	ArtifactFilename           string `json:"artifactFilename"`           // optional
}

type ExternalArtifact struct {
	TargetFolderRelativeToTask string `json:"targetFolderRelativeToTask"` //	optional
	Md5sum                     string `json:"md5sum"`                     // optional
	URL                        string `json:"url"`                        // optional
	Filename                   string `json:"filename"`                   //optional
	Filesize                   int64  `json:"filesize"`                   // long	// optional
	Name                       string `json:"name"`                       // optional
	IsArtifactList             bool   `json:"isArtifactList"`             //optional
}

type ExecutorData struct {
	SkipLogrotateAndCompress       bool                  `json:"skipLogrotateAndCompress"`       //	optional	If true, do not run logrotate or compress old log files
	LoggingExtraFields             map[string]string     `json:"loggingExtraFields"`             //Map[string,string]	optional
	EmbeddedArtifacts              []EmbeddedArtifact    `json:"embeddedArtifacts"`              //	optional	A list of the full content of any embedded artifacts
	S3Artifacts                    []S3Artifact          `json:"s3Artifacts"`                    //	optional	List of s3 artifacts for the executor to download
	SuccessfulExitCodes            []int                 `json:"successfulExitCodes"`            //	optional	Allowable exit codes for the task to be considered FINISHED instead of FAILED
	RunningSentinel                string                `json:"runningSentinel"`                // optional
	LogrotateFrequency             string                `json:"logrotateFrequency"`             // SingularityExecutorLogrotateFrequency	optional	Run logrotate this often. Can be HOURLY, DAILY, WEEKLY, MONTHLY
	MaxOpenFiles                   int                   `json:"maxOpenFiles"`                   // optional	Maximum number of open files the task process is allowed
	ExternalArtifacts              []ExternalArtifact    `json:"externalArtifacts"`              //	optional	A list of external artifacts for the executor to download
	User                           string                `json:"user"`                           // optional	Run the task process as this user
	PreserveTaskSandboxAfterFinish bool                  `json:"preserveTaskSandboxAfterFinish"` //optional	If true, do not delete files in the task sandbox after the task process has terminated
	ExtraCmdLineArgs               []string              `json:"extraCmdLineArgs"`               //	optional	Extra arguments in addition to any provided in the cmd field
	LoggingTag                     string                `json:"loggingTag"`                     //optional
	SigKillProcessesAfterMillis    int64                 `json:"sigKillProcessesAfterMillis"`    //long	optional	Send a sigkill to a process if it has not shut down this many millis after being sent a term signal
	MaxTaskThreads                 int                   `json:"maxTaskThreads"`                 // optional	Maximum number of threads a task is allowed to use
	S3ArtifactSignatures           []S3ArtifactSignature `json:"s3ArtifactSignatures"`           //	optional	A list of signatures use to verify downloaded s3artifacts
	Cmd                            string                `json:"cmd"`                            // required	Command for the custom executor to run
}

// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#model-SingularityMesosArtifact
type SingularityMesosArtifact struct {
	Cache      bool   `json:"cache,omitempty"`
	URI        string `json:"uri,omitempty"`
	Extract    bool   `json:"extract,omitempty"`
	Executable bool   `json:"executable,omitempty"`
}

// SingularityDeploy contains requird and optional parameter to configure
// a new and existing Singularity deploy.
type SingularityDeploy struct {
	CustomExecutorID                      string `json:"customExecutorId,omitempty"`
	SingularityDeployResources            `json:"resources,omitempty"`
	Uris                                  []SingularityMesosArtifact `json:"uris,omitempty"` //Array[SingularityMesosArtifact]	optional	List of URIs to download before executing the deploy command.
	ContainerInfo                         `json:"containerInfo"`
	Arguments                             []string                            `json:"arguments,omitempty"`
	TaskEnv                               interface{}                         `json:"taskEnv,omitempty"` // map[int]map[string]string //Map[int,Map[string,string]]	Map of environment variable overrides for specific task instances.
	AutoAdvanceDeploySteps                bool                                `json:"autoAdvanceDeploySteps,omitempty"`
	ServiceBasePath                       string                              `json:"serviceBasePath,omitempty"` // The base path for the API exposed by the deploy. Used in conjunction with the Load balancer API.
	CustomExecutorSource                  string                              `json:"customExecutorSource,omitempty"`
	Metadata                              map[string]string                   `json:"metadata,omitempty"`                 //ap of metadata key/value pairs associated with the deployment.
	TaskLabels                            map[int]map[string]string           `json:"taskLabels,omitempty"`               //Map[int,Map[string,string]]	optional	(Deprecated) Labels for specific tasks associated with this deploy, indexed by instance number
	MesosTaskLabels                       map[int][]SingularityMesosTaskLabel `json:"mesosTaskLabels,omitempty"`          // Map[int,List[SingularityMesosTaskLabel]] 	// optional	Labels for specific tasks associated with this deploy, indexed by instance number
	Labels                                map[string]string                   `json:"labels,omitempty"`                   // Map[string,string]	optional	Labels for all tasks associated with this deploy
	User                                  string                              `json:"user,omitempty"`                     //optional	Run tasks as this user
	RequestID                             string                              `json:"requestId"`                          // required	Singularity Request Id which is associated with this deploy.
	DeployStepWaitTimeMs                  int                                 `json:"deployStepWaitTimeMs,omitempty"`     // optional	wait this long between deploy steps
	SkipHealthchecksOnDeploy              bool                                `json:"skipHealthchecksOnDeploy,omitempty"` //optional	Allows skipping of health checks when deploying.
	MesosLabels                           *[]SingularityMesosTaskLabel        `json:"mesosLabels,omitempty"`              //Array[SingularityMesosTaskLabel]	optional	Labels for all tasks associated with this deploy
	Command                               string                              `json:"command,omitempty"`                  //optional	Command to execute for this deployment.
	*ExecutorData                         `json:"executorData,omitempty"`     //	optional	Executor specific information
	Shell                                 bool                                `json:"shell,omitempty"`                                 //optional	Override the shell property on the mesos task
	Timestamp                             int64                               `json:"timestamp,omitempty"`                             //long	optional	Deploy timestamp.
	DeployInstanceCountPerStep            int                                 `json:"deployInstanceCountPerStep,omitempty"`            //	optional	deploy this many instances at a time
	ConsiderHealthyAfterRunningForSeconds int64                               `json:"considerHealthyAfterRunningForSeconds,omitempty"` //	optional	Number of seconds that a service must be healthy to consider the deployment to be successful.
	MaxTaskRetries                        int                                 `json:"maxTaskRetries,omitempty"`                        // optional	allowed at most this many failed tasks to be retried before failing the deploy
	*SingularityRunNowRequest             `json:"runImmediately,omitempty"`   // optional	Settings used to run this deploy immediately
	CustomExecutorCmd                     string                              `json:"customExecutorCmd,omitempty"` // optional	Custom Mesos executor
	Env                                   map[string]string                   `json:"env,omitempty"`               //	optional	Map of environment variable definitions.
	// SingularityDeployResources            `json:"customExecutorResources"`    // com.hubspot.mesos.Resources	optional	Resources to allocate for custom mesos executor
	Version                    string `json:"version,omitempty"`                    //optional	Deploy version
	ID                         string `json:"id"`                                   //required	Singularity deploy id.
	DeployHealthTimeoutSeconds int64  `json:"deployHealthTimeoutSeconds,omitempty"` //optional	Number of seconds that Singularity waits for this service to become healthy (for it to download artifacts, start running, and optionally pass health
}

// SingularityDeployWithLB contains requird and optional parameter to configure
// a new and existing Singularity deploy.
type SingularityDeployWithLB struct {
	CustomExecutorID           string `json:"customExecutorId"`
	SingularityDeployResources `json:"resources"`
	Uris                       []string `json:"uris"` //Array[SingularityMesosArtifact]	optional	List of URIs to download before executing the deploy command.
	ContainerInfo              `json:"containerInfo"`
	// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#-set List of domains to host this service on, for use with the load balancer api
	LoadBalancerDomains                   []string `json:"loadBalancerDomains"` //Set
	HealthcheckOptions                    `json:"healthcheck"`
	Arguments                             []string          `json:"arguments"`
	TaskEnv                               interface{}       `json:"taskEnv"` // map[int]map[string]string //Map[int,Map[string,string]]	Map of environment variable overrides for specific task instances.
	AutoAdvanceDeploySteps                bool              `json:"autoAdvanceDeploySteps"`
	ServiceBasePath                       string            `json:"serviceBasePath"` // The base path for the API exposed by the deploy. Used in conjunction with the Load balancer API.
	CustomExecutorSource                  string            `json:"customExecutorSource"`
	Metadata                              map[string]string `json:"metadata"`                  //ap of metadata key/value pairs associated with the deployment.
	HealthcheckMaxRetries                 int               `json:"healthcheckMaxRetries"`     //optional	Maximum number of times to retry an individual healthcheck before failing the deploy.
	HealthcheckTimeoutSeconds             int64             `json:"healthcheckTimeoutSeconds"` //optional	Single healthcheck HTTP timeout in seconds.
	HealthcheckProtocol                   `json:"healthcheckProtocol"`
	TaskLabels                            map[int]map[string]string           `json:"taskLabels"`                        //Map[int,Map[string,string]]	optional	(Deprecated) Labels for specific tasks associated with this deploy, indexed by instance number
	HealthcheckPortIndex                  int                                 `json:"healthcheckPortIndex"`              //optional	Perform healthcheck on this dynamically allocated port (e.g. 0 for first port), defaults to first port
	HealthcheckMaxTotalTimeoutSeconds     int64                               `json:"healthcheckMaxTotalTimeoutSeconds"` //optional	Maximum amount of time to wait before failing a deploy for healthchecks to pass.
	LoadBalancerServiceIDOverride         string                              `json:"loadBalancerServiceIdOverride"`     //optional	Name of load balancer Service ID to use instead of the Request ID
	MesosTaskLabels                       map[int][]SingularityMesosTaskLabel `json:"mesosTaskLabels"`                   // Map[int,List[SingularityMesosTaskLabel]] 	// optional	Labels for specific tasks associated with this deploy, indexed by instance number
	Labels                                map[string]string                   `json:"labels"`                            // Map[string,string]	optional	Labels for all tasks associated with this deploy
	HealthcheckURI                        string                              `json:"healthcheckUri"`                    //optional	Deployment Healthcheck URI, if specified will be called after TASK_RUNNING.
	User                                  string                              `json:"user"`                              //optional	Run tasks as this user
	RequestID                             string                              `json:"requestId"`                         // required	Singularity Request Id which is associated with this deploy.
	LoadBalancerGroups                    interface{}                         `json:"loadBalancerGroups"`                // Set	optional	List of load balancer groups associated with this deployment.
	DeployStepWaitTimeMs                  int                                 `json:"deployStepWaitTimeMs"`              // optional	wait this long between deploy steps
	SkipHealthchecksOnDeploy              bool                                `json:"skipHealthchecksOnDeploy"`          //optional	Allows skipping of health checks when deploying.
	MesosLabels                           []SingularityMesosTaskLabel         `json:"mesosLabels"`                       //Array[SingularityMesosTaskLabel]	optional	Labels for all tasks associated with this deploy
	HealthcheckIntervalSeconds            int64                               `json:"healthcheckIntervalSeconds"`        //long	optional	Time to wait after a failed healthcheck to try again in seconds.
	Command                               string                              `json:"command"`                           //optional	Command to execute for this deployment.
	ExecutorData                          `json:"executorData"`               //	optional	Executor specific information
	LoadBalancerAdditionalRoutes          []string                            `json:"loadBalancerAdditionsRoutes"`           // optional	Additional routes besides serviceBasePath used by this service
	Shell                                 bool                                `json:"shell"`                                 //optional	Override the shell property on the mesos task
	Timestamp                             int64                               `json:"timestamp"`                             //long	optional	Deploy timestamp.
	DeployInstanceCountPerStep            int                                 `json:"deployInstanceCountPerStep"`            //	optional	deploy this many instances at a time
	ConsiderHealthyAfterRunningForSeconds int64                               `json:"considerHealthyAfterRunningForSeconds"` //	optional	Number of seconds that a service must be healthy to consider the deployment to be successful.
	LoadBalancerOptions                   map[string]interface{}              `json:"loadBalancerOptions"`                   // Map[string,Object]	optional	Map (Key/Value) of options for the load balancer.
	MaxTaskRetries                        int                                 `json:"maxTaskRetries"`                        // optional	allowed at most this many failed tasks to be retried before failing the deploy
	SingularityRunNowRequest              `json:"runImmediately"`             // optional	Settings used to run this deploy immediately
	LoadBalancerPortIndex                 int                                 `json:"loadBalancerPortIndex"`     // optional	Send this port to the load balancer api (e.g. 0 for first port), defaults to first port
	LoadBalancerTemplate                  string                              `json:"loadBalancerTemplate"`      // optional	Name of load balancer template to use if not using the default template
	CustomExecutorCmd                     string                              `json:"customExecutorCmd"`         // optional	Custom Mesos executor
	Env                                   map[string]string                   `json:"env"`                       //	optional	Map of environment variable definitions.
	LoadBalancerUpstreamGroup             string                              `json:"loadBalancerUpstreamGroup"` //optional	Group name to tag all upstreams with in load balancer
	// SingularityDeployResources            `json:"customExecutorResources"`    // com.hubspot.mesos.Resources	optional	Resources to allocate for custom mesos executor
	Version                    string `json:"version"`                    //optional	Deploy version
	ID                         string `json:"id"`                         //required	Singularity deploy id.
	DeployHealthTimeoutSeconds int64  `json:"deployHealthTimeoutSeconds"` //optional	Number of seconds that Singularity waits for this service to become healthy (for it to download artifacts, start running, and optionally pass health
}

type SingularityRunNowRequest struct {
	SingularityDeployResources `json:"resources,omitempty"` // optional	Override the resources from the active deploy for this run
	RunID                      string                       `json:"runId,omitempty"`            // optional	An id to associate with this request which will be associated with the corresponding launched tasks
	SkipHealthchecks           bool                         `json:"skipHealthchecks,omitempty"` // 	optional	If set to true, healthchecks will be skipped for this task run
	CommandLineArgs            []string                     `json:"commandLineArgs,omitempty"`  //	optional	Command line arguments to be passed to the task
	Message                    string                       `json:"message,omitempty"`          //optional	A message to show to users about why this action was taken
	RunAt                      int64                        `json:"runAt,omitEmpty"`            //long	optional	Schedule this task to run at a specified time
}

// SingularityExpiringPause contains information of a existing
// Singularity request.
type SingularityExpiringPause struct {
	User                                string      `json:"user"`
	RequestID                           string      `json:"requestId"`
	StartMillis                         int64       `json:"startMillis"`
	ActionID                            string      `json:"actionId"`
	SingularityExpiringAPIRequestObject interface{} `json:"expiringAPIRequestObject"`
}

// SingularityExpiringBounce contains information of a existing
// Singularity request.
type SingularityExpiringBounce struct {
	User                     string      `json:"user"`
	RequestID                string      `json:"requestId"`
	StartMillis              int64       `json:"startMillis"`
	DeployID                 string      `json:"deployId"`
	ActionID                 string      `json:"actionId"`
	ExpiringAPIRequestObject interface{} `json:"expiringAPIRequestObject"`
}

// SingularityDeployProgress contains deploy progress of a existing
// Singularity request.
type SingularityDeployProgress struct {
	AutoAdvanceDeploySteps     bool        `json:"autoAdvanceDeploySteps"`
	StepComplete               bool        `json:"stepComplete"`
	DeployStepWaitTimeMs       int64       `json:"deployStepWaitTimeMs"`
	Timestamp                  int64       `json:"timestamp"`
	DeployInstanceCountPerStep int         `json:"deployInstanceCountPerStep"`
	FailedDeployTasks          interface{} `json:"failedDeployTasks"` //Set	optional
	CurrentActiveInstances     int         `json:"currentActiveInstances"`
	TargetActiveInstances      int         `json:"targetActiveInstances"`
}

// SingularityLoadBalancerRequestID have loadbalancer information of a
// Singularity request.
type SingularityLoadBalancerRequestID struct {
	//optional	Allowable values: ADD, REMOVE, DEPLOY, DELETE
	RequestType   string `json:"requestType"`
	AttemptNumber int    `json:"attemptNumber"`
	ID            string `json:"id"`
}

// SingularityLoadBalancerUpdate contains parameters required to update a
// Singularity request's loadbalancer.
type SingularityLoadBalancerUpdate struct {
	// Allowable values: UNKNOWN, FAILED, WAITING, SUCCESS, CANCELING, CANCELED, INVALID_REQUEST_NOOP
	LoadBalancerState                string `json:"loadBalancerState"`
	SingularityLoadBalancerRequestID `json:"loadBalancerRequestId"`
	URI                              string `json:"uri"`
	// Allowable values: PRE_ENQUEUE, ENQUEUE, CHECK_STATE, CANCEL, DELETE
	Method    string `json:"method"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// SingularityDeployMarker holds information of a Singularity deploy.
type SingularityDeployMarker struct {
	User      string `json:"user"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	DeployID  string `json:"deployId"`
}

// SingularityDeployState holds information of a existing Singularity
// deploy.
type SingularityDeployState struct {
	// Allowable values: SUCCEEDED, FAILED_INTERNAL_STATE, CANCELING, WAITING, OVERDUE, FAILED, CANCELED
	CurrentDeployState            string `json:"currentDeployState"`
	SingularityRequest            `json:"updatedRequest"`
	SingularityDeployProgress     `json:"deployProgress"`
	SingularityLoadBalancerUpdate `json:"lastLoadBalancerUpdate"`
	SingularityDeployMarker       `json:"deployMarker"`
}

type SingularityExpiringAPIRequestObject struct {
	ActionID         string `json:"actionId"`
	DurationMillis   int64  `json:"durationMillis"`
	Instances        int64  `json:"instances"`
	Message          string `json:"message"`
	SkipHealthchecks bool   `json:"skipHealthchecks"`
}

// SingularityExpiringScale holds information of a expiring scale Singularity
// deploy.
type SingularityExpiringScale struct {
	RevertToInstances                   int    `json:"revertToInstances"`
	User                                string `json:"user"`
	RequestID                           string `json:"requestId"`
	Bounce                              bool   `json:"bounce"`
	StartMillis                         int64  `json:"startMillis"`
	ActionID                            string `json:"actionId"`
	DurationMillis                      int64  `json:"durationMillis"`
	SingularityExpiringAPIRequestObject `json:"expiringAPIRequestObject"`
}

// SingularityPendingDeploy holds information of a pending Singularity
// deploy.
type SingularityPendingDeploy struct {
	CurrentDeployState            string `json:"currentDeployState"` // Allowable values: SUCCEEDED, FAILED_INTERNAL_STATE, CANCELING, WAITING, OVERDUE, FAILED, CANCELED
	SingularityRequest            `json:"updatedRequest"`
	SingularityDeployProgress     `json:"deployProgress"`
	SingularityLoadBalancerUpdate `json:"lastLoadBalancerUpdate"`
	SingularityDeployMarker       `json:"deployMarker"`
}

// SingularityRequestParent contains request of a Singularity deploy.
type SingularityRequestParent struct {
	SingularityExpiringSkipHealthchecks `json:"expiringSkipHealthchecks"`
	PendingDeploy                       struct {
		CustomExecutorID           string `json:"customExecutorId"`
		SingularityDeployResources `json:"resources"`
		Uris                       interface{} `json:"uris"` //Array[SingularityMesosArtifact]	optional	List of URIs to download before executing the deploy command.
		ContainerInfo              `json:"containerInfo"`
		// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#-set List of domains to host this service on, for use with the load balancer api
		LoadBalancerDomains    interface{} `json:"loadBalancerDomains"` //Set
		HealthcheckOptions     `json:"healthcheck"`
		Arguments              []string    `json:"arguments"`
		TaskEnv                interface{} `json:"taskEnv"` //Map[int,Map[string,string]]	Map of environment variable overrides for specific task instances.
		AutoAdvanceDeploySteps bool        `json:"autoAdvanceDeploySteps"`
		ID                     string      `json:"id"`
		Command                string      `json:"command"`
	} `json:"pendingDeploy"`
	ActiveDeploy struct {
		CustomExecutorID           string `json:"customExecutorId"`
		SingularityDeployResources `json:"resources"`
		Uris                       interface{} `json:"uris"` //Array[SingularityMesosArtifact]	optional	List of URIs to download before executing the deploy command.
		ContainerInfo              `json:"containerInfo"`
		// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#-set List of domains to host this service on, for use with the load balancer api
		LoadBalancerDomains    interface{} `json:"loadBalancerDomains"` //Set
		HealthcheckOptions     `json:"healthcheck"`
		Arguments              []string    `json:"arguments"`
		TaskEnv                interface{} `json:"taskEnv"` //Map[int,Map[string,string]]	Map of environment variable overrides for specific task instances.
		AutoAdvanceDeploySteps bool        `json:"autoAdvanceDeploySteps"`
		ID                     string      `json:"id"`
		Command                string      `json:"command"`
	} `json:"activeDeploy"`
	SingularityExpiringPause  `json:"expiringPause"`
	SingularityExpiringBounce `json:"expiringBounce"`
	SingularityRequest        `json:"request"`
	SingularityPendingDeploy  `json:"pendingDeployState"`
	SingularityExpiringScale  `json:"expiringScale"`
	RequestDeployState        struct {
		ActiveDeploy struct {
			DeployID  string `json:"deployId"`
			RequestID string `json:"requestId"`
			Timestamp int64  `json:"timestamp"`
		} `json:"activeDeploy"`
		RequestID string `json:"requestId"`
	} `json:"requestDeployState"`
	State string `json:"state"`
}

// SingularityDeleteRequest contains HTTP body for a Delete Singularity Request. Please see below URL for
// more information.
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#delete-apirequestsrequestrequestid
// https://github.com/HubSpot/Singularity/blob/master/Docs/reference/api.md#-singularitydeleterequestrequest
type SingularityDeleteRequest struct {
	DeleteFromLoadBalancer bool   `json:"deleteFromLoadBalancer"` //optional	Should the service associated with the request be removed from the load balancer
	Message                string `json:"message"`                //optional	A message to show to users about why this action was taken
	ActionID               string `json:"actionId"`               //An id to associate with this action for metadata purposes
}

type SingularityDeployRequest struct {
	UnpauseOnSuccessfulDeploy bool                    //optional	If deploy is successful, also unpause the request
	SingularityDeploy         `json:"deploy"`         // required	The Singularity deploy object, containing all the required details about the Deploy
	*SingularityRequest       `json:"updatedRequest"` // optional	use this request data for this deploy, and update the request on successful deploy
	Message                   string                  //optional	A message to show users about this deploy (metadata)
}
