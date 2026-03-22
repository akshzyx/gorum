package user

import "time"

type PublicProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Bio       string `json:"bio" validate:"max=200"`
	AvatarURL string `json:"avatar_url" validate:"url"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" validate:"email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"min=8"`
}
