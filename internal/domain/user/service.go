package user

import (
	"context"

	"github.com/akshzyx/gorum/internal/util"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByEmail is mostly internal use (auth), returns raw repo result
func (s *Service) GetByEmail(ctx context.Context, email string) (User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (s *Service) GetPublicProfile(
	ctx context.Context,
	username string,
) (PublicProfileResponse, error) {
	profile, err := s.repo.GetPublicByUsername(ctx, username)
	if err != nil {
		return PublicProfileResponse{}, ErrNotFound
	}
	return profile, nil
}

func (s *Service) UpdateEmail(
	ctx context.Context,
	userID string,
	req UpdateEmailRequest,
) error {
	if req.Email == "" {
		return ErrForbidden
	}

	return s.repo.UpdateEmail(ctx, userID, req.Email)
}

func (s *Service) UpdatePassword(
	ctx context.Context,
	userID string,
	req UpdatePasswordRequest,
) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return ErrNotFound
	}

	if !util.CheckPasswordHash(req.CurrentPassword, u.PasswordHash) {
		return ErrInvalidPassword
	}

	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, hash)
}
