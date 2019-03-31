package main

import (
	"flag"
	"marathon-alerts/alert/slack"
	"marathon-alerts/marathon"

	"github.com/sirupsen/logrus"
)

var (
	marathonUrl string
)

func main() {
	flag.StringVar(&marathonUrl, "marathonUrl", "http://127.0.0.1:8080", "Uri for marathon")
	flag.Parse()

	logrus.SetLevel(logrus.TraceLevel)

	eventHandler := marathon.Handler{Url: marathonUrl, Alertor: slack.Slack{}}
	eventHandler.Init()
}
