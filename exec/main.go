package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bperreault-va/extra-life-slack-notifier/discord"
	"github.com/bperreault-va/extra-life-slack-notifier/extralife"
	"github.com/bperreault-va/extra-life-slack-notifier/slack"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	config, err := readConfig()
	if err != nil {
		config, err = createConfig()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = writeConfig(config)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var input string
	for input != "4" {
		input, err = printMenu(config)
		if err != nil {
			fmt.Println(err.Error())
			_, err = reader.ReadString('\n')
			return
		}
		switch input {
		case "1":
			err = updateTeamID(config)
		case "2":
			err = updateSlackWebhookURL(config)
		case "3":
			err = updateDiscordWebhookURL(config)
		case "":
			return
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		err = writeConfig(config)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var slackService slack.Service
	var discordService discord.Service
	if config.SlackWebhookURL != "" {
		slackService = slack.New(config.SlackWebhookURL)
	}
	if config.DiscordWebhookURL != "" {
		discordService = discord.New(config.DiscordWebhookURL)
	}
	extraLifeService := extralife.New(config.TeamID, slackService, discordService)

	var slackTeam extralife.Team
	fmt.Println("Testing Extra Life Team ID...")
	slackTeam, err = extraLifeService.GetTeam()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to find team with ID %s.", config.TeamID))
		fmt.Println(err.Error())
		fmt.Println("Press enter to exit.")
		_, err = reader.ReadString('\n')
		return
	}
	fmt.Println(fmt.Sprintf("Slack Team name: %s", slackTeam.Name))

	if slackService != nil {
		fmt.Println("Testing Slack Incoming Webhook URL...")
		err = slackService.SendTestSlacktivity()
		if err != nil {
			fmt.Println(fmt.Sprintf("Testing Slack Incoming Webhook URL failed with the following status: %s", err.Error()))
			_, err = reader.ReadString('\n')
			return
		}
	}

	if discordService != nil {
		fmt.Println("Testing Discord Incoming Webhook URL...")
		err = discordService.SendTestMessage()
		if err != nil {
			fmt.Println(fmt.Sprintf("Testing Discord Webhook URL failed with the following status: %s", err.Error()))
			_, err = reader.ReadString('\n')
			return
		}
	}

	fmt.Println("Starting server...")
	extraLifeService.PollExtraLife()
}

type config struct {
	TeamID            string `json:"teamID"`
	SlackWebhookURL   string `json:"slackWebhookURL"`
	DiscordWebhookURL string `json:"discordWebhookURL"`
}

func writeConfig(c *config) error {
	f, err := os.Create("elconfig.json")
	if err != nil {
		return err
	}
	defer f.Close()
	confJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = f.Write(confJSON)
	if err != nil {
		return err
	}
	return nil
}

func readConfig() (*config, error) {
	var c *config
	data, err := ioutil.ReadFile("./elconfig.json")
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func createConfig() (c *config, err error) {
	conf := &config{}
	if err = updateTeamID(conf); err != nil {
		return conf, err
	}
	if err = updateSlackWebhookURL(conf); err != nil {
		return conf, err
	}
	if err = updateDiscordWebhookURL(conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func updateTeamID(c *config) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Extra Life Team ID:")
	teamID, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	c.TeamID = strings.TrimSpace(teamID)
	return nil
}

func updateSlackWebhookURL(c *config) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Slack Incoming Webhook URL:")
	slackWebhookURL, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	c.SlackWebhookURL = strings.TrimSpace(slackWebhookURL)
	return nil
}

func updateDiscordWebhookURL(c *config) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Discord Webhook URL:")
	discordWebhookURL, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	c.DiscordWebhookURL = strings.TrimSpace(discordWebhookURL)
	return nil
}

func printMenu(c *config) (string, error) {
	fmt.Println(fmt.Sprintf("\nExtra Life Team ID:\t%s", c.TeamID))
	fmt.Println(fmt.Sprintf("Slack Webhook URL:\t%s", c.SlackWebhookURL))
	fmt.Println(fmt.Sprintf("Discord Webhook URL:\t%s", c.DiscordWebhookURL))
	fmt.Println("1) Update Team ID")
	fmt.Println("2) Update Slack Webhook URL")
	fmt.Println("3) Update Discord Webhook URL")
	fmt.Println("4) Start server")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}
