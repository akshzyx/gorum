package user

import (
	"context"
	"errors"

	"github.com/akshzyx/gorum/internal/util"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) UpdateProfile(context context.Context, param any, req UpdateProfileRequest) any {
	panic("unimplemented")
}

func (s *Service) GetByEmail(ctx context.Context, email string) (User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetPublicProfile(
	ctx context.Context,
	username string,
) (PublicProfileResponse, error) {
	return s.repo.GetPublicByUsername(ctx, username)
}

// func (s *Service) UpdateProfile(
// 	ctx context.Context,
// 	userID string,
// 	req UpdateProfileRequest,
// ) error {
// 	return s.repo.UpdateProfile(ctx, userID, req.Bio, req.AvatarURL)
// }

func (s *Service) UpdateEmail(
	ctx context.Context,
	userID string,
	req UpdateEmailRequest,
) error {
	return s.repo.UpdateEmail(ctx, userID, req.Email)
}

func (s *Service) UpdatePassword(
	ctx context.Context,
	userID string,
	req UpdatePasswordRequest,
) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !util.CheckPasswordHash(req.CurrentPassword, u.PasswordHash) {
		return errors.New("invalid password")
	}

	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, hash)
}
