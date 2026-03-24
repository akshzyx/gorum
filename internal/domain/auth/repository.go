package auth

import (
	"context"
)

type Repository interface {
	// USER LOOKUPS
	GetUserByEmail(ctx context.Context, email string) (AuthUser, error)
	GetUserByID(ctx context.Context, id string) (AuthUser, error)

	// USER CREATION (Signup)
	// This MUST run inside a transaction
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (AuthUser, error)

	// EMAIL VERIFICATION TOKEN
	CreateVerificationToken(ctx context.Context, token VerificationToken) error
	GetVerificationToken(ctx context.Context, token string) (VerificationToken, error)
	MarkVerificationTokenUsed(ctx context.Context, token string) error

	// Activate user
	ActivateUser(ctx context.Context, userID string) error
}
