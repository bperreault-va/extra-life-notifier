package main

import (
	"fmt"
	"github.com/bperreault-va/extra-life-slack-notifier/discord"
	"github.com/bperreault-va/extra-life-slack-notifier/extralife"
	"github.com/bperreault-va/extra-life-slack-notifier/slack"
	"net/http"
)

func main() {
	fmt.Println("Starting HTTP server...")
	handleHTTP()

	fmt.Println("Starting Slack server...")

	slackService := slack.New("https://hooks.slack.com/services/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")
	discordService := discord.New("https://discordapp.com/api/webhooks/632072354436218880/tAFR4RkKu5KSZ62Wu83Oh0MV1HkzDUsPvS6TsR-50fl59UECcx7Bwrkl8Mw8jNC6O9KN")
	s := extralife.New("44339", slackService, discordService)
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
