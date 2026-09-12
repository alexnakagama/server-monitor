package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  error
	}{
		{
			name:     "empty username is required",
			username: "",
			wantErr:  errors_custom.ErrUsernameRequired,
		},
		{
			name:     "username shorter than 5 characters is too short",
			username: "abcd",
			wantErr:  errors_custom.ErrUsernameTooShort,
		},
		{
			name:     "username of exactly 5 characters is valid",
			username: "abcde",
			wantErr:  nil,
		},
		{
			name:     "username of 30 characters is valid",
			username: strings.Repeat("a", 30),
			wantErr:  nil,
		},
		{
			name:     "username longer than 30 characters is too long",
			username: strings.Repeat("a", 31),
			wantErr:  errors_custom.ErrUsernameTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateUsername(tt.username)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ValidateUsername(%q) error = %v, want %v", tt.username, gotErr, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{
			name:    "empty email is required",
			email:   "",
			wantErr: errors_custom.ErrEmailRequired,
		},
		{
			name:    "email longer than 255 characters is too long",
			email:   strings.Repeat("a", 250) + "@b.com",
			wantErr: errors_custom.ErrEmailTooLong,
		},
		{
			name:    "malformed email is invalid",
			email:   "not-an-email",
			wantErr: errors_custom.ErrInvalidEmail,
		},
		{
			name:    "email with display name is invalid",
			email:   "John Doe <john@example.com>",
			wantErr: errors_custom.ErrInvalidEmail,
		},
		{
			name:    "email missing domain is invalid",
			email:   "user@",
			wantErr: errors_custom.ErrInvalidEmail,
		},
		{
			name:    "simple valid email is valid",
			email:   "user@example.com",
			wantErr: nil,
		},
		{
			name:    "email of exactly 255 characters is valid",
			email:   strings.Repeat("a", 249) + "@b.com",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateEmail(tt.email)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ValidateEmail(%q) error = %v, want %v", tt.email, gotErr, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "empty password is required",
			password: "",
			wantErr:  errors_custom.ErrPasswordRequired,
		},
		{
			name:     "password shorter than 6 characters is too short",
			password: "abcde",
			wantErr:  errors_custom.ErrPasswordTooShort,
		},
		{
			name:     "password of exactly 6 characters is valid",
			password: "abcdef",
			wantErr:  nil,
		},
		{
			name:     "password of 80 characters is valid",
			password: strings.Repeat("a", 80),
			wantErr:  nil,
		},
		{
			name:     "password longer than 80 characters is too long",
			password: strings.Repeat("a", 81),
			wantErr:  errors_custom.ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidatePassword(tt.password)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ValidatePassword(%q) error = %v, want %v", tt.password, gotErr, tt.wantErr)
			}
		})
	}
}
