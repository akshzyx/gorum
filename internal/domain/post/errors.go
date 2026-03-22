package post

import "errors"

var (
	ErrPostNotFound = errors.New("post not found")
	ErrUnauthorised = errors.New("unauthorised")
)
