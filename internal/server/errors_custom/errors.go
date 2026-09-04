package errors_custom

import "errors"

var (
	ErrServerNotFound = errors.New("server not found")
	ErrUserNotFound   = errors.New("user not found")

	ErrNameRequired     = errors.New("name is required")
	ErrHostnameRequired = errors.New("hostname is required")
	ErrOSRequired       = errors.New("os is required")
)
