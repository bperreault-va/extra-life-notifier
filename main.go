package main

import (
	"fmt"
	"github.com/branson-perreault/extra-life-notifier/discord"
	"github.com/branson-perreault/extra-life-notifier/extralife"
	"github.com/branson-perreault/extra-life-notifier/slack"
	"net/http"
)

func main() {
	fmt.Println("Starting HTTP server...")
	handleHTTP()

	teamID := "CHANGE ME"
	slackWebhookURL := "CHANGE ME OR LEAVE BLANK"
	discordWebhookURL := "CHANGE ME OR LEAVE BLANK"

	slackService := slack.New(slackWebhookURL)
	discordService := discord.New(discordWebhookURL)
	s := extralife.New(teamID, slackService, discordService)

	err := testMessaging(slackService, discordService)
	if err != nil {
		fmt.Println("Shutting down")
		return
	}
	fmt.Println("Polling extra-life.org...")
	s.PollExtraLife()
}

func handleHTTP() {
	// Hello world
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello Extra Life!"))
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	// Health check
	http.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err.Error())
		}
	}()
}

func testMessaging(slackService slack.Service, discordService discord.Service) error {
	fmt.Println("  Testing Slack server...")
	if slackService.IsConfigured() {
		err := slackService.SendMessage("  Restarting Slack server...")
		if err != nil {
			fmt.Println(fmt.Sprintf("    Slack test failure: %s", err.Error()))
			return err
		}
		fmt.Println("    Slack test success!")
	} else {
		fmt.Println("    Slack is not configured. Skipping.")
	}
	fmt.Println("  Testing Discord server...")
	if discordService.IsConfigured() {
		err := discordService.SendMessage("  Restarting Discord server...")
		if err != nil {
			fmt.Println(fmt.Sprintf("    Discord test failure: %s", err.Error()))
			return err
		}
		fmt.Println("    Discord test success!")
	} else {
		fmt.Println("    Discord is not configured. Skipping.")
	}
	return nil
}
