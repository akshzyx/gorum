package postgres

import (
	"context"

	db "github.com/akshzyx/gorum/internal/db/sqlc/generated"
	"github.com/akshzyx/gorum/internal/domain/follow"
	"github.com/akshzyx/gorum/internal/domain/user"
)

type FollowRepository struct {
	q *db.Queries
}

func NewFollowRepository(q *db.Queries) follow.Repository {
	return &FollowRepository{
		q: q,
	}
}

func (r *FollowRepository) CreateFollow(ctx context.Context, followerID, followingID string) error {
	return r.q.CreateFollow(ctx, db.CreateFollowParams{
		FollowerID:  followerID,
		FollowingID: followingID,
	})
}

func (r *FollowRepository) DeleteFollow(ctx context.Context, followerID, followingID string) error {
	return r.q.DeleteFollow(ctx, db.DeleteFollowParams{
		FollowerID:  followerID,
		FollowingID: followingID,
	})
}

func (r *FollowRepository) GetFollowers(ctx context.Context, userID string) ([]user.User, error) {
	rows, err := r.q.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	var users []user.User
	for _, u := range rows {
		users = append(users, user.User{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Time,
		})
	}

	return users, nil
}

func (r *FollowRepository) GetFollowing(ctx context.Context, userID string) ([]user.User, error) {
	rows, err := r.q.GetFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	var users []user.User
	for _, u := range rows {
		users = append(users, user.User{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Time,
		})
	}

	return users, nil
}
