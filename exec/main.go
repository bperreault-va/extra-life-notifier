package main

import (
	"bufio"
	"fmt"
	"github.com/bperreault-va/extra-life-slack-notifier/slack"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Extra Life Team ID:")
	teamID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	teamID = strings.TrimSpace(teamID)

	fmt.Println("Slack Incoming Webhook URL:")
	incomingWebhookURL, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	incomingWebhookURL = strings.TrimSpace(incomingWebhookURL)

	fmt.Println(fmt.Sprintf("Proceed using Team ID: %s Webhook URL: %s? (Y/n)", teamID, incomingWebhookURL))
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		_, err = reader.ReadString('\n')
		return
	}
	if answer != "\n" && strings.Contains(answer, "n") {
		return
	}

	fmt.Println("Testing Extra Life Team ID...")
	slackService := slack.New(teamID, incomingWebhookURL)
	team, err := slackService.GetTeam()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to find team with ID %s.", teamID))
		fmt.Println(err.Error())
		fmt.Println("Press enter to exit.")
		_, err = reader.ReadString('\n')
		return
	}

	fmt.Println("Testing Slack Incoming Webhook URL...")
	err = slackService.SendTestSlacktivity()
	if err != nil {
		fmt.Println(fmt.Sprintf("Testing Slack Incoming Webhook URL failed with the following status: %s", err.Error()))
		_, err = reader.ReadString('\n')
		return
	}

	fmt.Println(fmt.Sprintf("Team name: %s", team.Name))
	fmt.Println("Starting server...")
	slackService.PollExtraLife()
}
