package errors_custom

import "errors"

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrEmailRequired    = errors.New("email is required")

	ErrServerNotFound = errors.New("server not found")
	ErrUserNotFound   = errors.New("user not found")

	ErrNameRequired     = errors.New("name is required")
	ErrHostnameRequired = errors.New("hostname is required")
	ErrOSRequired       = errors.New("os is required")

	ErrHostnameAlreadyExists = errors.New("hostname already exists")
)
