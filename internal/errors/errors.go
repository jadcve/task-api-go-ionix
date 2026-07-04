package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("inactive user")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCurrentPass = errors.New("invalid current password")
	ErrUnauthorized       = errors.New("unauthorized")
)
