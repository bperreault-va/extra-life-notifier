package discord

type Service interface {
	SendTestMessage() error
	SendMessage(message string) error
}
