package extralife

type Service interface {
	PollExtraLife()
	GetTeam() (Team, error)
	GetParticipants() ([]Participant, error)
	GetRecentDonations() ([]Donation, error)
}
