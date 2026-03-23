package post

import "context"

type Repository interface {
	// Posts
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id string) (*Post, error)
	Delete(ctx context.Context, postID string) error
	ListLatest(ctx context.Context, limit int32) ([]*Post, error)

	// Replies
	GetPostForReply(ctx context.Context, id string) (*Post, error)
	CreateReply(ctx context.Context, post *Post) error
	ListReplies(ctx context.Context, postID string) ([]*Post, error)
	GetThread(ctx context.Context, rootID string) ([]*Post, error)
}
