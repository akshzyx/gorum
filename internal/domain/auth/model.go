package auth

import "time"


type VerificationToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

type PasswordResetToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

type AuthUser struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	IsVerified   bool
}

type CreateUserTxParams struct {
	ID                string
	Username          string
	Email             string
	PasswordHash      string
	VerificationToken VerificationToken
}
