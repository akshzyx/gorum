package user

import "time"

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	IsVerified   bool
	CreatedAt    time.Time
}
