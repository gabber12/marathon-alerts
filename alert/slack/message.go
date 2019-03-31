package slack

const (
	GREEN  = "#4caf50"
	RED    = "#f44336"
	YELLOW = "#ffeb3b"
)

type SlackMessage struct {
	Blocks      []SlackBlock      `json:"blocks"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Text      string          `json:"text,omitempty"`
	Accessory *SlackAccessory `json:"accessory,omitempty"`
	Color     *string         `json:"color"`
}

type SlackBlockElement struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackBlock struct {
	Type      string             `json:"type"`
	Text      *SlackBlockElement `json:"text,omitempty"`
	Accessory *SlackAccessory    `json:"accessory,omitempty"`
}

type SlackAccessory struct {
	Type string            `json:"type"`
	Text SlackBlockElement `json:"text"`
	Url  string            `json:"url"`
}
