package main

import (
	"fmt"
	"github.com/bperreault-va/extra-life-slack-notifier/slack"
	"net/http"
)

func main() {
	fmt.Println("Starting HTTP server...")
	handleHTTP()

	fmt.Println("Starting Slack server...")
	s := slack.New("44172", "https://hooks.slack.com/services/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")
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
