package enums

type ProfileStatus int64

const (
	ProfileStatusFresh ProfileStatus = 1 << iota
	ProfileStatusContestantFinished
	ProfileStatusOrganizerFinished
)
