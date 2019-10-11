package slack

type Service interface {
	PollExtraLife()
	GetTeam() (Team, error)
	GetParticipants() ([]Participant, error)
	GetRecentDonations() ([]Donation, error)
	SendDonationSlacktivity(team Team, donation Donation, participant Participant) error
	SendParticipantSlacktivity(participant Participant) error
	SendTestSlacktivity() error
}
