package user

import "time"

type PublicProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Bio       *string   `json:"bio,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type MeResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
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
