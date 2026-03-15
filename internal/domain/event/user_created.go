package event

import (
	"time"
)

const UserCreatedEventType = "user.created"

type UserCreatedEvent struct {
	userID     string
	name       string
	email      string
	occurredAt time.Time
}

func NewUserCreatedEvent(userID, name, email string) *UserCreatedEvent {
	return &UserCreatedEvent{
		userID:     userID,
		name:       name,
		email:      email,
		occurredAt: time.Now(),
	}
}

func (e *UserCreatedEvent) EventType() string {
	return UserCreatedEventType
}

func (e *UserCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *UserCreatedEvent) AggregateID() string {
	return e.userID
}

func (e *UserCreatedEvent) UserID() string {
	return e.userID
}

func (e *UserCreatedEvent) Name() string {
	return e.name
}

func (e *UserCreatedEvent) Email() string {
	return e.email
}

