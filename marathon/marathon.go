package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type App struct {
	Data AppData `json:"app"`
}
type AppData struct {
	Id        string            `json:"id"`
	Instances int32             `json:"instances"`
	Staged    int32             `json:"tasksStaged"`
	Running   int32             `json:"tasksRunning"`
	Healthy   int32             `json:"tasksHealthy"`
	UnHealthy int32             `json:"tasksUnHealthy"`
	Labels    map[string]string `json:"labels"`
	Failure   *TaskFailure      `json:"lastTaskFailure"`
}
type TaskFailure struct {
	Message string `json:"message"`
}

type Client struct {
	Url            string
	MarathonLeader string
	MesosLeader    string
	FrameworkId    string
}

type Info struct {
	FrameworkId string         `json:"frameworkId"`
	LeaderHost  string         `json:"leader"`
	Config      MarathonConfig `json:"marathon_config"`
}

type MarathonConfig struct {
	MesosLeader string `json:"mesos_leader_ui_url"`
}

func (c *Client) GetAppSpec(appId string) (*App, error) {
	resp, err := http.Get(c.Url + "/v2/apps/" + appId)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var app App
	if err = json.Unmarshal(body, &app); err == nil {
		return &app, err
	}
	return nil, fmt.Errorf("Error getting App Spec: %v", err)
}

func (c *Client) GetInfo() (*Info, error) {
	resp, err := http.Get(c.Url + "/v2/info")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var app Info
	if err = json.Unmarshal(body, &app); err == nil {
		c.MarathonLeader = app.LeaderHost
		c.MesosLeader = app.Config.MesosLeader
		c.FrameworkId = app.FrameworkId
		return &app, err
	}

	return nil, fmt.Errorf("Error getting App Spec: %v", err)
}
