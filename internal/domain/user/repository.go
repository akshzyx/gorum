package user

import "context"

type Repository interface {
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetPublicByUsername(ctx context.Context, username string) (PublicProfileResponse, error)

	// UpdateProfile(ctx context.Context, userID, bio, avatarURL string) error
	UpdateEmail(ctx context.Context, userID, email string) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error

	GetByID(ctx context.Context, id string) (User, error)
}
