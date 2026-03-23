package post

import (
	"context"
	"strings"

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

func (s *Service) Create(ctx context.Context, userID string, req *CreatePostRequest) (string, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return "", ErrInvalidContent
	}

	p := &Post{
		ID:      ulid.Make().String(),
		UserID:  userID,
		Content: content,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return "", err
	}

	return p.ID, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

// Delete checks ownership in the service layer before deleting the post.
func (s *Service) Delete(ctx context.Context, postID, userID string) error {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return ErrPostNotFound
	}

	// Only the owner of the post should be allowed to delete it.
	if post.UserID != userID {
		return ErrForbidden
	}

	if err := s.repo.Delete(ctx, postID); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListLatest(ctx context.Context, limit int32) ([]*Post, error) {
	return s.repo.ListLatest(ctx, limit)
}

func (s *Service) Reply(ctx context.Context, userID string, parentID string, req *CreateReplyRequest) (string, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return "", ErrInvalidContent
	}

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
		Content:      content,
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
