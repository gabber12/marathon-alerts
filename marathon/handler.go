package marathon

import (
	"encoding/json"
	"fmt"
	"marathon-alerts/alert"
	"net/url"

	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	MarathonClient *Client
	Url            string
	Alertor        alert.Alertor
}

func (h *Handler) Init() error {
	if h.MarathonClient == nil {
		h.MarathonClient = &Client{Url: h.Url}
	}

	info, err := h.MarathonClient.GetInfo()
	if err != nil {
		logrus.Errorf("Error getting version info from marathon %v", err)
		return err
	}
	logrus.Infof("Marathon %s Leader: %s", info.FrameworkId, info.LeaderHost)
	err = h.subcribeToEventBus()
	if err != nil {
		logrus.Errorf("Unable to Subscribe to marathon event bus")
	}
	return err
}

func (h *Handler) Unmarshal(payload []byte) (interface{}, error) {
	var event Event
	err := json.Unmarshal(payload, &event)
	if err != nil {
		return nil, err
	}
	if event.Type == STATUS_UPDATE_EVENT {
		var res StatusUpdateEvent
		err = json.Unmarshal(payload, &res)
		return res, err
	} else if event.Type == HEALTH_CHECK_FAILURE_EVENT {
		var res HealthCheckFailedEvent
		err = json.Unmarshal(payload, &res)
		return res, err
	}
	return nil, err
}

func (h *Handler) Handle(payload []byte) error {
	event, err := h.Unmarshal(payload)
	if err != nil {
		logrus.Errorf("Error Marshalling payload %v", err)
		return err
	}
	switch ev := event.(type) {
	case StatusUpdateEvent:
		h.handleStatus(ev)
	case HealthCheckFailedEvent:

	}
	return nil
}

func (h Handler) shouldHandle(app App) bool {
	alertLabel := "ALERT_ON"
	for k, _ := range app.Data.Labels {
		if k == alertLabel {
			return true
		}
	}
	return false
}

func (h *Handler) subcribeToEventBus() error {
	client := sse.NewClient(h.MarathonClient.Url + "/v2/events")
	logrus.Tracef("Subscribing to Event bus...")

	return client.Subscribe("events", func(msg *sse.Event) {
		logrus.Tracef("Event Received %+v", string(msg.Data))
		go h.Handle(msg.Data)
	})
}

func (h *Handler) buildAlert(app App, event StatusUpdateEvent) alert.AppAlert {
	id := url.QueryEscape(app.Data.Id)
	sandboxUrl := fmt.Sprintf("%s/#/agents/%s/frameworks/%s/executors/%s/tasks/%s/browse", h.MarathonClient.MesosLeader, *event.SlaveId, h.MarathonClient.FrameworkId, *event.TaskId, *event.TaskId)
	appUrl := fmt.Sprintf("%s/ui/#/apps/%s", h.Url, id)
	return alert.AppAlert{
		App:               app.Data.Id,
		Level:             alert.WARNING,
		RequiredInstances: app.Data.Instances,
		HealthyInstances:  app.Data.Healthy,
		Host:              *event.TaskHost,
		MarathonAppURL:    appUrl,
		FailureReason:     app.Data.Failure.Message,
		MesosSandboxUrl:   sandboxUrl,
		Name:              *event.TaskStatus,
	}
}
func (h *Handler) handleStatus(event StatusUpdateEvent) error {
	app, err := h.MarathonClient.GetAppSpec(*event.AppId)
	if err != nil {
		logrus.Errorf("Error Handling Event %v", err)
		return nil
	}
	if h.shouldHandle(*app) {
		logrus.Infof("Handling App Event [%s] -> Event Type [%s] -> Status [%s]", app.Data.Id, event.Type, *event.TaskStatus)
		if *event.TaskStatus == TS_KILLED || *event.TaskStatus == TS_FAILED {
			alert := h.buildAlert(*app, event)
			h.Alertor.Alert(alert)
		}
	} else {
		logrus.Infof("Alerting Off for app: [%s]. Ignoring.", app.Data.Id)
	}
	return nil
}
