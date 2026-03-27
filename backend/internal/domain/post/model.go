package post

import "time"

type Post struct {
	ID           string
	UserID       string
	Username     string
	Content      string
	Likes        int64
	Liked        bool
	ParentPostID *string
	RootPostID   *string
	CreatedAt    time.Time
}
