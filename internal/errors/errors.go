package errors

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInactiveUser         = errors.New("inactive user")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCurrentPass   = errors.New("invalid current password")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidUserRole      = errors.New("invalid user role")
	ErrTaskNotFound         = errors.New("task not found")
	ErrInvalidTaskStatus    = errors.New("invalid task status")
	ErrInvalidDueDate       = errors.New("invalid due date")
	ErrAssignedUserNotFound = errors.New("assigned user not found")
)
