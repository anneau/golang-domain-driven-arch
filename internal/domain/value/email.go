package value

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	value string
}

func NewEmail(s string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	if !emailRegexp.MatchString(normalized) {
		return Email{}, fmt.Errorf("invalid email format: %s", s)
	}
	return Email{value: normalized}, nil
}

func EmailFrom(s string) Email {
	return Email{value: s}
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
