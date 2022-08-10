package discord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (s *service) IsConfigured() bool {
	return s.webhookURL != ""
}

func (s *service) SendTestMessage() error {
	data := url.Values{}
	data.Add("content", "Starting server...")
	resp, err := http.PostForm(s.webhookURL, data)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(resp.Status, "2") {
		return fmt.Errorf("%v", resp.Status)
	}
	return nil
}

func (s *service) SendMessage(message string) error {
	data := url.Values{}
	data.Add("content", message)
	resp, err := http.PostForm(s.webhookURL, data)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(resp.Status, "2") {
		return fmt.Errorf("%v", resp.Status)
	}
	return nil
}

type Webhook struct {
	Token string `json:"token"`
}

func (s *service) GetToken() (string, error) {
	resp, err := http.PostForm(s.webhookURL, url.Values{})
	if err != nil {
		return "", err
	}
	if !strings.Contains(resp.Status, "200") {
		fmt.Println(fmt.Sprintf("%v", resp.Status))
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var webhook Webhook
	err = json.Unmarshal([]byte(response), &webhook)
	if err != nil {
		return "", err
	}
	return webhook.Token, nil
}
