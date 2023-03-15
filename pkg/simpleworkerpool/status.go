package simpleworkerpool

const (
	MinNonState StatusType = iota
	RunningStatus
	IdleStatus
	BusyStatus
	ClosedStatus
	StoppedStatus

	MaxNonState
)
