package user

import "errors"

var (
	ErrNotFound        = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrForbidden       = errors.New("forbidden")
)
