package slack

type Service interface {
	IsConfigured() bool
	SendMessage(message string) error
	SendTestMessage() error
}
