package goevent

import "context"

type (
	EventName     string
	EventPriority int
)

const (
	HighPriority EventPriority = 3
	MedPriority  EventPriority = 2
	LowPriority  EventPriority = 1
)

type Event interface {
	Name() EventName
	Priority() EventPriority
}

type Listener interface {
	Handle(ctx context.Context, event Event) error
}
