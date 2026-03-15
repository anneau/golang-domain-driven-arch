package entity

import (
	"fmt"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/event"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/value"
)

type User struct {
	id                value.UserID
	name              string
	email             value.Email
	uncommittedEvents []event.Event
}

func NewUser(id value.UserID, name string, email value.Email) (*User, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	u := &User{id: id, name: name, email: email}
	u.uncommittedEvents = append(u.uncommittedEvents, event.NewUserCreatedEvent(
		id.String(), name, email.String(),
	))
	return u, nil
}

func (u *User) UncommittedEvents() []event.Event {
	return u.uncommittedEvents
}

func (u *User) ClearEvents() {
	u.uncommittedEvents = nil
}

func ReconstructUser(id value.UserID, name string, email value.Email) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
	}
}

func (u *User) ID() value.UserID   { return u.id }
func (u *User) Name() string       { return u.name }
func (u *User) Email() value.Email { return u.email }

func (u *User) ChangeName(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	u.name = name
	return nil
}

func (u *User) ChangeEmail(email value.Email) {
	u.email = email
}

func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return u.id.Equals(other.id)
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("user name must not be empty")
	}
	if len([]rune(name)) > 100 {
		return fmt.Errorf("user name must be 100 characters or less")
	}
	return nil
}
