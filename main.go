package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GetParticipantsURL                   = "https://extra-life.org/api/teams/44172/participants"
	GetDonationsURL                      = "https://extra-life.org/api/teams/44172/donations?limit=5"
	SlackWebhookURL                      = "https://hooks.slack.com/services/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"
	SlacktivityTemplate                  = "%s just received a $%.2f donation from %s!"
	SlacktivityAdditionalMessageTemplate = "\n> %s"
	ParticipantSlacktivityTemplate       = "%s joined the team!"
	TimeLayout                           = "2006-01-02T15:04:05"
	WaitDuration                         = 60 * time.Second
)

func main() {
	// Hello world
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Extra Life!")
	})
	// Health check
	http.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	var participants []Participant
	var err error
	for {
		participants, err = GetParticipants()
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, participant := range participants {
			created, err := time.Parse(TimeLayout, strings.Split(participant.Created, ".")[0])
			if err != nil {
				fmt.Println(err.Error())
			}
			if created.After(time.Now().UTC().Add(-WaitDuration)) {
				err := SendParticipantSlacktivity(participant)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		donations, err := GetRecentDonations()
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, donation := range donations {
			created, err := time.Parse(TimeLayout, strings.Split(donation.Created, ".")[0])
			if err != nil {
				fmt.Println(err.Error())
			}
			if created.Before(time.Now().UTC().Add(-WaitDuration)) {
				continue
			}
			for _, p := range participants {
				if p.ParticipantID == donation.ParticipantID {
					err = SendDonationSlacktivity(donation, p)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
		time.Sleep(WaitDuration)
	}
}

type Participant struct {
	ParticipantID int64  `json:"participantID"`
	DisplayName   string `json:"displayName"`
	Created       string `json:"createdDateUTC"`
}

func GetParticipants() ([]Participant, error) {
	resp, err := http.Get(GetParticipantsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var participants []Participant
	err = json.Unmarshal([]byte(response), &participants)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

type Donation struct {
	DisplayName   string  `json:"displayName"`
	Message       string  `json:"message"`
	ParticipantID int64   `json:"participantID"`
	Amount        float64 `json:"amount"`
	Created       string  `json:"createdDateUTC"`
}

func GetRecentDonations() ([]Donation, error) {
	resp, err := http.Get(GetDonationsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var donations []Donation
	err = json.Unmarshal([]byte(response), &donations)
	if err != nil {
		return nil, err
	}
	return donations, nil
}

type Slacktivity struct {
	Text string `json:"text"`
}

func SendDonationSlacktivity(donation Donation, participant Participant) error {
	message := fmt.Sprintf(SlacktivityTemplate, participant.DisplayName, donation.Amount, donation.DisplayName)
	if donation.Message != "" {
		message += fmt.Sprintf(SlacktivityAdditionalMessageTemplate, donation.Message)
	}
	fmt.Println(message)
	slacktivity := Slacktivity{Text: message}
	payload, err := json.Marshal(slacktivity)
	if err != nil {
		return err
	}
	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(SlackWebhookURL, data)
	if err != nil {
		return err
	}
	if resp.Status != "200" {
		fmt.Println(resp)
	}
	return nil
}

func SendParticipantSlacktivity(participant Participant) error {
	message := fmt.Sprintf(ParticipantSlacktivityTemplate, participant.DisplayName)
	fmt.Println(message)
	slacktivity := Slacktivity{Text: message}
	payload, err := json.Marshal(slacktivity)
	if err != nil {
		return err
	}
	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(SlackWebhookURL, data)
	if err != nil {
		return err
	}
	if resp.Status != "200" {
		fmt.Println(resp)
	}
	return nil
}
