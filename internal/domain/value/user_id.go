package value

import (
	"crypto/rand"
	"fmt"
)

type UserID struct {
	value string
}

func NewUserID() (UserID, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return UserID{}, fmt.Errorf("failed to generate user id: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	id := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return UserID{value: id}, nil
}

func UserIDFrom(s string) (UserID, error) {
	if s == "" {
		return UserID{}, fmt.Errorf("user id must not be empty")
	}
	return UserID{value: s}, nil
}

func (id UserID) String() string {
	return id.value
}

func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}
