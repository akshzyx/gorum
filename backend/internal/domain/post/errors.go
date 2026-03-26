package post

import "errors"

var (
	ErrPostNotFound   = errors.New("post not found")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidContent = errors.New("content cannot be empty")
)
