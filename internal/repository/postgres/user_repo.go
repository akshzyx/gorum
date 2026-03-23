package postgres

import (
	"context"
	"errors"

	db "github.com/akshzyx/gorum/internal/db/sqlc/generated"
	"github.com/akshzyx/gorum/internal/domain/user"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) Create(ctx context.Context, id, username, email, passwordHash string) error {
	params := db.CreateUserParams{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
	return r.q.CreateUser(ctx, params)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, err
	}

	return user.User{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsVerified:   row.IsVerified,
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, err
	}

	return user.User{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsVerified:   row.IsVerified,
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

// public profile only returns safe fields
func (r *UserRepository) GetPublicByUsername(
	ctx context.Context,
	username string,
) (user.PublicProfileResponse, error) {
	row, err := r.q.GetPublicProfileByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.PublicProfileResponse{}, user.ErrNotFound
		}
		return user.PublicProfileResponse{}, err
	}

	return user.PublicProfileResponse{
		ID:        row.ID,
		Username:  row.Username,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (user.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, err
	}

	return user.User{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsVerified:   row.IsVerified,
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) UpdateEmail(
	ctx context.Context,
	userID, email string,
) error {
	return r.q.UpdateUserEmail(ctx, db.UpdateUserEmailParams{
		ID:    userID,
		Email: email,
	})
}

func (r *UserRepository) UpdatePassword(
	ctx context.Context,
	userID, hash string,
) error {
	return r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: hash,
	})
}
