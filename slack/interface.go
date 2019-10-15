package slack

type Service interface {
	SendSlackMessage(message string) error
	SendTestSlacktivity() error
}
