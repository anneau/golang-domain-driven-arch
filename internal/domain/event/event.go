package event

import "time"

type Event interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() string
}
