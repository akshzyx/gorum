package postgres

import (
	"context"
	"errors"

	db "github.com/akshzyx/gorum/internal/db/sqlc/generated"
	"github.com/akshzyx/gorum/internal/domain/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	q    *db.Queries
	pool *pgxpool.Pool
}

func NewAuthRepo(q *db.Queries, pool *pgxpool.Pool) auth.Repository {
	return &AuthRepo{
		q:    q,
		pool: pool,
	}
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (auth.AuthUser, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.AuthUser{}, auth.ErrUserNotFound
		}
		return auth.AuthUser{}, err
	}

	return auth.AuthUser{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsVerified:   row.IsVerified,
	}, nil
}

func (r *AuthRepo) GetUserByID(ctx context.Context, id string) (auth.AuthUser, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.AuthUser{}, auth.ErrUserNotFound
		}
		return auth.AuthUser{}, err
	}

	return auth.AuthUser{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsVerified:   row.IsVerified,
	}, nil
}

func (r *AuthRepo) CreateUserTx(ctx context.Context, arg auth.CreateUserTxParams) (auth.AuthUser, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return auth.AuthUser{}, err
	}

	qtx := r.q.WithTx(tx)

	// 1. Insert user
	if err := qtx.CreateUser(ctx, db.CreateUserParams{
		ID:           arg.ID,
		Username:     arg.Username,
		Email:        arg.Email,
		PasswordHash: arg.PasswordHash,
	}); err != nil {
		tx.Rollback(ctx)
		return auth.AuthUser{}, err
	}

	// 2. Insert verification token
	vt := arg.VerificationToken

	if err := qtx.CreateVerificationToken(ctx, db.CreateVerificationTokenParams{
		Token:     vt.Token,
		UserID:    vt.UserID,
		ExpiresAt: pgtype.Timestamptz{Time: vt.ExpiresAt, Valid: true},
	}); err != nil {
		tx.Rollback(ctx)
		return auth.AuthUser{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return auth.AuthUser{}, err
	}

	return auth.AuthUser{
		ID:       arg.ID,
		Username: arg.Username,
		Email:    arg.Email,
	}, nil
}

func (r *AuthRepo) CreateVerificationToken(ctx context.Context, t auth.VerificationToken) error {
	return r.q.CreateVerificationToken(ctx, db.CreateVerificationTokenParams{
		Token:     t.Token,
		UserID:    t.UserID,
		ExpiresAt: pgtype.Timestamptz{Time: t.ExpiresAt, Valid: true},
	})
}

func (r *AuthRepo) GetVerificationToken(ctx context.Context, token string) (auth.VerificationToken, error) {
	row, err := r.q.GetVerificationToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.VerificationToken{}, auth.ErrTokenNotFound
		}
		return auth.VerificationToken{}, err
	}

	return auth.VerificationToken{
		Token:     row.Token,
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt.Time,
		Used:      row.Used,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *AuthRepo) MarkVerificationTokenUsed(ctx context.Context, token string) error {
	return r.q.MarkVerificationTokenUsed(ctx, token)
}

func (r *AuthRepo) ActivateUser(ctx context.Context, userID string) error {
	return r.q.ActivateUser(ctx, userID)
}
