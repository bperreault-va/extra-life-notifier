package discord

type Service interface {
	IsConfigured() bool
	SendTestMessage() error
	SendMessage(message string) error
}
