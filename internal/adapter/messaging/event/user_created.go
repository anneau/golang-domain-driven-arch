package event

import "time"

type UserCreatedEventDTO struct {
	EventType  string    `json:"event_type"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	OccurredAt time.Time `json:"occurred_at"`
}
