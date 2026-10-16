package auth

import (
	"errors"
)

var (
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password must have at least 8 characters")
	ErrPasswordTooLong    = errors.New("password must have at most 72 bytes")
	ErrEmailInUse         = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
)
