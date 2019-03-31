package marathon

const (
	STATUS_UPDATE_EVENT        = "status_update_event"
	HEALTH_CHECK_FAILURE_EVENT = "failed_health_check_event"
)

const (
	TS_KILLED = "TASK_KILLED"
	TS_FAILED = "TASK_FAILED"
)

type HealthCheckFailedEvent struct {
	*Event
	AppId  *string `json:"appId,omitempty"`
	TaskId *string `json:"taskId"`
}

type StatusUpdateEvent struct {
	*Event
	*StatusUpdateData
}

type Event struct {
	Type      string `json:"eventType"`
	Timestamp string `json:"timestamp"`
}

type StatusUpdateData struct {
	AppId      *string `json:"appId,omitempty"`
	TaskStatus *string `json:"taskStatus,omitempty"`
	TaskHost   *string `json:"host,omitempty"`
	SlaveId    *string `json:"slaveId,omitempty"`
	TaskId     *string `json:"taskId,omitempty"`
}
