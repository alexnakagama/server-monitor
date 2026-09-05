package model

import (
	"time"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *User) Validate() {
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors_custom.ErrUsernameRequired
	}

	if len(username) < 5 {
		return errors_custom.ErrUsernameTooShort
	}

	if len(username) > 30 {
		return errors_custom.ErrUsernameTooLong
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors_custom.ErrEmailRequired
	}

	return nil
}
