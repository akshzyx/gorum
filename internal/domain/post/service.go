package post

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userId string, req *CreatePostRequest) (string, error) {
	p := &Post{
		ID:      ulid.Make().String(),
		UserID:  userId,
		Content: req.Content,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return "", err
	}

	return p.ID, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) DeleteByOwner(ctx context.Context, postID, userID string) error {
	return s.repo.DeleteByOwner(ctx, postID, userID)
}

func (s *Service) ListLatest(ctx context.Context, limit int32) ([]*Post, error) {
	return s.repo.ListLatest(ctx, limit)
}

func (s *Service) Reply(ctx context.Context, userID string, parentID string, req *CreateReplyRequest) (string, error) {
	parent, err := s.repo.GetPostForReply(ctx, parentID)
	if err != nil {
		return "", ErrPostNotFound
	}

	rootID := parent.ID
	if parent.RootPostID != nil {
		rootID = *parent.RootPostID
	}

	reply := &Post{
		ID:           ulid.Make().String(),
		UserID:       userID,
		Content:      req.Content,
		ParentPostID: &parentID,
		RootPostID:   &rootID,
	}

	if err := s.repo.CreateReply(ctx, reply); err != nil {
		return "", err
	}

	return reply.ID, nil
}

func (s *Service) ListReplies(ctx context.Context, postID string) ([]*Post, error) {
	return s.repo.ListReplies(ctx, postID)
}

func (s *Service) GetThread(ctx context.Context, rootID string) ([]*Post, error) {
	return s.repo.GetThread(ctx, rootID)
}
