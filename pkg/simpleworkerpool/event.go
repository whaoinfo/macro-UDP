package simpleworkerpool

type EventType int

const (
	MinNonEvent EventType = iota
	CloseEvent
	ExecuteEvent
	ExceptionEvent
	CompleteEvent

	MaxNonEvent
)
