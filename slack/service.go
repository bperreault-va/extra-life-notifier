package slack

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
	GetTeamURL                                   = "https://extra-life.org/api/teams/%s"
	GetParticipantsURL                           = "https://extra-life.org/api/teams/%s/participants"
	GetDonationsURL                              = "https://extra-life.org/api/teams/%s/donations?limit=5"
	DonationSlacktivityTemplate                  = "%s just received a $%.2f donation from %s!"
	DonationSlacktivityAdditionalMessageTemplate = "\n> %s"
	DonationSlacktivityTotalTemplate             = "\nNew team total: $%.2f"
	ParticipantSlacktivityTemplate               = "%s joined the team!"
	TimeLayout                                   = "2006-01-02T15:04:05"
	WaitDuration                                 = 60 * time.Second
)

type slackService struct {
	teamID          string
	slackWebhookURL string
}

func New(teamID string, slackWebhookURL string) Service {
	return &slackService{
		teamID:          teamID,
		slackWebhookURL: slackWebhookURL,
	}
}

func (s *slackService) PollExtraLife() {
	var participants []Participant
	var err error
	for {
		participants, err = s.GetParticipants()
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, participant := range participants {
			created, err := time.Parse(TimeLayout, strings.Split(participant.Created, ".")[0])
			if err != nil {
				fmt.Println(err.Error())
			}
			if created.After(time.Now().UTC().Add(-WaitDuration)) {
				err := s.SendParticipantSlacktivity(participant)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		donations, err := s.GetRecentDonations()
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
					team, err := s.GetTeam()
					if err != nil {
						fmt.Println(err.Error())
					}

					err = s.SendDonationSlacktivity(team, donation, p)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
		time.Sleep(WaitDuration)
	}
}

type Team struct {
	FundraisingGoal float64 `json:"fundraisingGoal"`
	SumDonations    float64 `json:"sumDonations"`
	Name            string  `json:"name"`
}

func (s *slackService) GetTeam() (Team, error) {
	resp, err := http.Get(fmt.Sprintf(GetTeamURL, s.teamID))
	if err != nil {
		return Team{}, err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Team{}, err
	}

	var team Team
	err = json.Unmarshal([]byte(response), &team)
	if err != nil {
		return Team{}, err
	}
	return team, nil
}

type Participant struct {
	ParticipantID int64  `json:"participantID"`
	DisplayName   string `json:"displayName"`
	Created       string `json:"createdDateUTC"`
}

func (s *slackService) GetParticipants() ([]Participant, error) {
	resp, err := http.Get(fmt.Sprintf(GetParticipantsURL, s.teamID))
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

func (s *slackService) GetRecentDonations() ([]Donation, error) {
	resp, err := http.Get(fmt.Sprintf(GetDonationsURL, s.teamID))
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

func (s *slackService) SendDonationSlacktivity(team Team, donation Donation, participant Participant) error {
	message := fmt.Sprintf(DonationSlacktivityTemplate, participant.DisplayName, donation.Amount, donation.DisplayName)
	if donation.Message != "" {
		message += fmt.Sprintf(DonationSlacktivityAdditionalMessageTemplate, donation.Message)
	}
	message += fmt.Sprintf(DonationSlacktivityTotalTemplate, team.SumDonations)

	fmt.Println(message)
	slacktivity := Slacktivity{Text: message}
	payload, err := json.Marshal(slacktivity)
	if err != nil {
		return err
	}
	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(s.slackWebhookURL, data)
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Status, "200") {
		fmt.Println(fmt.Sprintf("%s", resp.Status))
	}
	return nil
}

func (s *slackService) SendParticipantSlacktivity(participant Participant) error {
	message := fmt.Sprintf(ParticipantSlacktivityTemplate, participant.DisplayName)
	fmt.Println(message)
	slacktivity := Slacktivity{Text: message}
	payload, err := json.Marshal(slacktivity)
	if err != nil {
		return err
	}
	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(s.slackWebhookURL, data)
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Status, "200") {
		fmt.Println(fmt.Sprintf("%v", resp.Status))
	}
	return nil
}

func (s *slackService) SendTestSlacktivity() error {
	payload, err := json.Marshal(Slacktivity{Text: "Starting server..."})
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Add("payload", string(payload))
	resp, err := http.PostForm(s.slackWebhookURL, data)
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Status, "200") {
		return fmt.Errorf("%v", resp.Status)
	}
	return nil
}
