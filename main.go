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

	fmt.Println("Starting Slack server...")

	slackService := slack.New("")
	discordService := discord.New("")
	s := extralife.New("", slackService, discordService)
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
