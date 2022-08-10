package extralife

import (
	"encoding/json"
	"fmt"
	"github.com/branson-perreault/extra-life-notifier/discord"
	"github.com/branson-perreault/extra-life-notifier/slack"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	GetTeamURL                        = "https://www.extra-life.org/api/teams/%s"
	GetParticipantsURL                = "https://www.extra-life.org/api/teams/%s/participants"
	GetDonationsURL                   = "https://www.extra-life.org/api/teams/%s/donations?limit=5"
	DonationMessageTemplate           = "%s just received a $%.2f USD donation from %s! 💸"
	DonationAdditionalMessageTemplate = "\n> %s"
	DonationTeamTotalTemplate         = "\nNew team total: $%.2f USD 💸"
	ParticipantMessageTemplate        = "%s joined the team! 🎉"
	TimeLayout                        = "2006-01-02T15:04:05"
	WaitDuration                      = 60 * time.Second
)

type extralifeService struct {
	teamID         string
	slackService   slack.Service
	discordService discord.Service
}

func New(teamID string, slack slack.Service, discord discord.Service) Service {
	return &extralifeService{
		teamID:         teamID,
		slackService:   slack,
		discordService: discord,
	}
}

func (s *extralifeService) PollExtraLife() {
	var participants []Participant
	var err error
	for {
		participants, err = s.GetParticipants()
		if err != nil {
			fmt.Println(fmt.Sprintf("Error getting participants: %s", err.Error()))
		}
		for _, participant := range participants {
			created, err := time.Parse(TimeLayout, strings.Split(participant.Created, ".")[0])
			if err != nil {
				fmt.Println(fmt.Sprintf("Error parsing participant time: %s", err.Error()))
			}
			if created.After(time.Now().UTC().Add(-WaitDuration)) {
				err := s.SendParticipantMessage(participant)
				if err != nil {
					fmt.Println(fmt.Sprintf("Error sending participant message: %s", err.Error()))
				}
			}
		}
		donations, err := s.GetRecentDonations()
		if err != nil {
			fmt.Println(fmt.Sprintf("Error getting recent donations: %s", err.Error()))
		}
		for _, donation := range donations {
			created, err := time.Parse(TimeLayout, strings.Split(donation.Created, ".")[0])
			if err != nil {
				fmt.Println(fmt.Sprintf("Error parsing donation time: %s", err.Error()))
			}
			if created.Before(time.Now().UTC().Add(-WaitDuration)) {
				continue
			}
			for _, p := range participants {
				if p.ParticipantID == donation.ParticipantID {
					team, err := s.GetTeam()
					if err != nil {
						fmt.Println(fmt.Sprintf("Error getting team: %s", err.Error()))
					}

					err = s.SendDonationMessage(team, donation, p)
					if err != nil {
						fmt.Println(fmt.Sprintf("Error sending donation message: %s", err.Error()))
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

func (s *extralifeService) GetTeam() (Team, error) {
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

func (s *extralifeService) GetParticipants() ([]Participant, error) {
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
	err = json.Unmarshal(response, &participants)
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

func (s *extralifeService) GetRecentDonations() ([]Donation, error) {
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
	err = json.Unmarshal(response, &donations)
	if err != nil {
		return nil, err
	}
	return donations, nil
}

func (s *extralifeService) SendDonationMessage(team Team, donation Donation, participant Participant) error {
	message := fmt.Sprintf(DonationMessageTemplate, participant.DisplayName, donation.Amount, donation.DisplayName)
	if donation.Message != "" {
		message += fmt.Sprintf(DonationAdditionalMessageTemplate, donation.Message)
	}
	message += fmt.Sprintf(DonationTeamTotalTemplate, team.SumDonations)

	if s.slackService.IsConfigured() {
		err := s.slackService.SendMessage(message)
		if err != nil {
			return err
		}
	}

	if s.discordService.IsConfigured() {
		err := s.discordService.SendMessage(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *extralifeService) SendParticipantMessage(participant Participant) error {
	message := fmt.Sprintf(ParticipantMessageTemplate, participant.DisplayName)

	if s.slackService.IsConfigured() {
		err := s.slackService.SendMessage(message)
		if err != nil {
			return err
		}
	}
	if s.discordService.IsConfigured() {
		err := s.discordService.SendMessage(message)
		if err != nil {
			return err
		}
	}
	return nil
}
