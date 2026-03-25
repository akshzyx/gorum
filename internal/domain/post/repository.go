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
	CountReplies(ctx context.Context, postID string) (int64, error)
	GetThread(ctx context.Context, rootID string) ([]*Post, error)

	// Likes
	CreateLike(ctx context.Context, userID, postID string) error
	DeleteLike(ctx context.Context, userID, postID string) error
	CountLikes(ctx context.Context, postID string) (int64, error)
	HasUserLiked(ctx context.Context, userID, postID string) (bool, error)

	// Likes (batch, used for feeds/lists)
	GetLikesCountByPostIDs(ctx context.Context, postIDs []string) (map[string]int64, error)
	GetUserLikedPosts(ctx context.Context, userID string, postIDs []string) (map[string]bool, error)

	// User profile posts and replies
	GetPostsByUser(ctx context.Context, userID string, limit int32) ([]*Post, error)
	GetRepliesByUser(ctx context.Context, userID string, limit int32) ([]*Post, error)
}
