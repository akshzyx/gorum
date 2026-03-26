package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrTokenNotFound      = errors.New("verification token not found")
	ErrTokenExpired       = errors.New("verification token expired")
	ErrTokenUsed          = errors.New("verification token already used")
	ErrAlreadyVerified    = errors.New("email already verified")
)
