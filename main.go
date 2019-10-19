package main

import (
	"fmt"
	"github.com/bperreault-va/extra-life-notifier/discord"
	"github.com/bperreault-va/extra-life-notifier/extralife"
	"github.com/bperreault-va/extra-life-notifier/slack"
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

	fmt.Println("Starting Slack server...")
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
