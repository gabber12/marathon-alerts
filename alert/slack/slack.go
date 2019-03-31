package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"marathon-alerts/alert"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Slack struct {
}

func (s Slack) sendMessage(text []byte) error {

	resp, err := http.Post("https://hooks.slack.com/services/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", "application/json", bytes.NewBuffer(text))
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		logrus.Errorf("Slack Notification Status Code %d %s", resp.StatusCode, body)
		return fmt.Errorf("Problem Posting to slack channel. Response Code %d %s", resp.StatusCode, body)
	}
	return nil
}

func (s Slack) Alert(alert alert.AppAlert) error {
	block :=
		SlackBlock{
			Type: "section",
			Text: &SlackBlockElement{
				Type: "mrkdwn",
				Text: fmt.Sprintf("<%s|*%s*> :: *%s* on `%s`", alert.MarathonAppURL, alert.App, alert.Name, alert.Host),
			},
		}
	divider := SlackBlock{Type: "divider", Text: nil}
	alertTypeBlock := SlackBlock{
		Type: "section",
		Text: &SlackBlockElement{
			Type: "mrkdwn",
			Text: fmt.Sprintf("```%s```", alert.FailureReason),
		},
		Accessory: &SlackAccessory{
			Type: "button",
			Text: SlackBlockElement{
				Type: "plain_text",
				Text: "See Logs",
			},
			Url: alert.MesosSandboxUrl,
		},
	}

	color := RED
	message := SlackMessage{
		Blocks: []SlackBlock{
			block,
			alertTypeBlock,
			divider,
		},
		Attachments: []SlackAttachment{
			SlackAttachment{
				Text:  fmt.Sprintf("*Instance State*\nHealthy: *%d* \nRequired Instances: *%d*", alert.HealthyInstances, alert.RequiredInstances),
				Color: &color,
			},
		},
	}
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Error Marshling slack message %v", err)
	}
	logrus.Info("Slack Message %s", string(body))
	s.sendMessage(body)
	return nil
}
