package domain

import "errors"

// Sentinel errors. Services return these; handlers map them to HTTP status.
var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotFound           = errors.New("not found")
	ErrForbidden          = errors.New("forbidden")
)
