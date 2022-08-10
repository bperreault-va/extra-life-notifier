package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type service struct {
	webhookURL string
}

func New(webhookURL string) Service {
	return &service{
		webhookURL: webhookURL,
	}
}

type slackMessage struct {
	Text string `json:"text"`
}

func (s *service) IsConfigured() bool {
	return s.webhookURL != ""
}

func (s *service) SendMessage(message string) error {

	fmt.Println(message)
	slacktivity := slackMessage{Text: message}
	payload, err := json.Marshal(slacktivity)
	if err != nil {
		return err
	}
	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(s.webhookURL, data)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(resp.Status, "2") {
		fmt.Println(fmt.Sprintf("%v", resp.Status))
	}
	return nil
}

func (s *service) SendTestMessage() error {
	payload, err := json.Marshal(slackMessage{Text: "Starting server..."})
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(s.webhookURL, data)
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Status, "200") {
		return fmt.Errorf("%v", resp.Status)
	}
	return nil
}
