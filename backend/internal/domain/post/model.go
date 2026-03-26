package post

import "time"

type Post struct {
	ID           string
	UserID       string
	Username     string
	Content      string
	ParentPostID *string
	RootPostID   *string
	CreatedAt    time.Time
}
