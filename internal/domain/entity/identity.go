package entity

import "github.com/hkobori/golang-domain-driven-arch/internal/domain/value"

type Identity struct {
	email  value.Email
	userID value.UserID
}

func NewIdentity(userID value.UserID, email value.Email) *Identity {
	return &Identity{
		userID: userID,
		email:  email,
	}
}

func (i *Identity) UserID() value.UserID { return i.userID }
func (i *Identity) Email() value.Email   { return i.email }
