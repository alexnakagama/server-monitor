package errors_custom

import "errors"

var (
	ErrServerNotFound = errors.New("server not found")
	ErrUserNotFound   = errors.New("user not found")
)
