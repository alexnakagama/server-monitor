package errors_custom

import "errors"

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")

	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")

	ErrUsernameTooShort = errors.New("username is too short")
	ErrEmailTooShort    = errors.New("email is too short")
	ErrPasswordTooShort = errors.New("password is too short")

	ErrUsernameTooLong = errors.New("username is too long")
	ErrEmailTooLong    = errors.New("email is too long")
	ErrPasswordTooLong = errors.New("password is too long")

	ErrServerNotFound = errors.New("server not found")
	ErrUserNotFound   = errors.New("user not found")

	ErrNameRequired     = errors.New("name is required")
	ErrHostnameRequired = errors.New("hostname is required")
	ErrOSRequired       = errors.New("os is required")

	ErrHostnameAlreadyExists = errors.New("hostname already exists")
)
