package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type service struct {
	webhookUrl string
}

func New(webhookUrl string) Service {
	return &service{
		webhookUrl: webhookUrl,
	}
}

type slackMessage struct {
	Text string `json:"text"`
}

func (s *service) SendSlackMessage(message string) error {

	fmt.Println(message)
	slacktivity := slackMessage{Text: message}
	payload, err := json.Marshal(slacktivity)
	if err != nil {
		return err
	}
	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(s.webhookUrl, data)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(resp.Status, "2") {
		fmt.Println(fmt.Sprintf("%v", resp.Status))
	}
	return nil
}

func (s *service) SendTestSlacktivity() error {
	payload, err := json.Marshal(slackMessage{Text: "Starting server..."})
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(s.webhookUrl, data)
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Status, "200") {
		return fmt.Errorf("%v", resp.Status)
	}
	return nil
}
