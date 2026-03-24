package follow

import (
	"context"

	"github.com/akshzyx/gorum/internal/domain/user"
)

type Repository interface {
	CreateFollow(ctx context.Context, followerID, followingID string) error
	DeleteFollow(ctx context.Context, followerID, followingID string) error

	GetFollowers(ctx context.Context, userID string) ([]user.User, error)
	GetFollowing(ctx context.Context, userID string) ([]user.User, error)
}
