package follow

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) FollowUser(ctx context.Context, followerID, followingID string) error {
	if followerID == followingID {
		return ErrCannotFollowYourself
	}

	return s.repo.CreateFollow(ctx, followerID, followingID)
}

func (s *Service) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	return s.repo.DeleteFollow(ctx, followerID, followingID)
}

func (s *Service) GetFollowers(ctx context.Context, userID string) ([]UserPreview, error) {
	users, err := s.repo.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []UserPreview
	for _, u := range users {
		result = append(result, UserPreview{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		})
	}

	return result, nil
}

func (s *Service) GetFollowing(ctx context.Context, userID string) ([]UserPreview, error) {
	users, err := s.repo.GetFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []UserPreview
	for _, u := range users {
		result = append(result, UserPreview{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		})
	}

	return result, nil
}
