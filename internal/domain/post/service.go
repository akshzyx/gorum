package post

import (
	"context"
	"strings"
	"time"

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

func (s *Service) ListLatest(
	ctx context.Context,
	userID string,
	cursor *time.Time,
	limit int32,
) (*PaginatedPosts, error) {
	rows, err := s.repo.ListLatestWithCursor(ctx, cursor, limit)
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, &Post{
			ID:        r.ID,
			UserID:    r.UserID,
			Content:   r.Content,
			CreatedAt: r.CreatedAt.Time,
		})
	}

	// enrich posts (likes, liked, reply_count)
	enriched, err := s.EnrichPosts(ctx, userID, posts)
	if err != nil {
		return nil, err
	}

	// next cursor
	var nextCursor *time.Time
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		nextCursor = &last.CreatedAt
	}

	hasMore := len(posts) == int(limit)

	return &PaginatedPosts{
		Data:       enriched,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
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

func (s *Service) ListReplies(
	ctx context.Context,
	userID string,
	postID string,
	cursor *time.Time,
	limit int32,
) (*PaginatedPosts, error) {
	rows, err := s.repo.ListRepliesWithCursor(ctx, postID, cursor, limit)
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, 0, len(rows))
	for _, r := range rows {
		var parentID *string
		if r.ParentPostID.Valid {
			parentID = &r.ParentPostID.String
		}

		var rootID *string
		if r.RootPostID.Valid {
			rootID = &r.RootPostID.String
		}

		posts = append(posts, &Post{
			ID:           r.ID,
			UserID:       r.UserID,
			Content:      r.Content,
			ParentPostID: parentID,
			RootPostID:   rootID,
			CreatedAt:    r.CreatedAt.Time,
		})
	}

	enriched, err := s.EnrichPosts(ctx, userID, posts)
	if err != nil {
		return nil, err
	}

	var nextCursor *time.Time
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		nextCursor = &last.CreatedAt
	}

	hasMore := len(posts) == int(limit)

	return &PaginatedPosts{
		Data:       enriched,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *Service) GetRepliesCount(ctx context.Context, postID string) (int64, error) {
	return s.repo.CountReplies(ctx, postID)
}

func (s *Service) GetThread(ctx context.Context, rootID string) ([]*Post, error) {
	return s.repo.GetThread(ctx, rootID)
}

// Likes (single ops)

func (s *Service) LikePost(ctx context.Context, userID, postID string) error {
	// make sure post exists
	_, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return ErrPostNotFound
	}

	return s.repo.CreateLike(ctx, userID, postID)
}

func (s *Service) UnlikePost(ctx context.Context, userID, postID string) error {
	return s.repo.DeleteLike(ctx, userID, postID)
}

func (s *Service) GetLikesCount(ctx context.Context, postID string) (int64, error) {
	return s.repo.CountLikes(ctx, postID)
}

func (s *Service) HasUserLiked(ctx context.Context, userID, postID string) (bool, error) {
	return s.repo.HasUserLiked(ctx, userID, postID)
}

func (s *Service) EnrichPosts(ctx context.Context, userID string, posts []*Post) ([]map[string]interface{}, error) {
	postIDs := make([]string, 0, len(posts))
	for _, p := range posts {
		postIDs = append(postIDs, p.ID)
	}

	likesMap, err := s.repo.GetLikesCountByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	likedMap := map[string]bool{}
	if userID != "" {
		likedMap, err = s.repo.GetUserLikedPosts(ctx, userID, postIDs)
		if err != nil {
			return nil, err
		}
	}

	resp := make([]map[string]interface{}, 0, len(posts))

	for _, p := range posts {
		replyCount, _ := s.repo.CountReplies(ctx, p.ID)

		resp = append(resp, map[string]interface{}{
			"id":          p.ID,
			"user_id":     p.UserID,
			"content":     p.Content,
			"created_at":  p.CreatedAt,
			"likes":       likesMap[p.ID],
			"liked":       likedMap[p.ID],
			"reply_count": replyCount,
		})
	}

	return resp, nil
}

func (s *Service) GetUserPosts(
	ctx context.Context,
	userID string,
	cursor *time.Time,
	limit int32,
	viewerID string,
) (*PaginatedPosts, error) {
	rows, err := s.repo.GetPostsByUserWithCursor(ctx, userID, cursor, limit)
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, &Post{
			ID:        r.ID,
			UserID:    r.UserID,
			Content:   r.Content,
			CreatedAt: r.CreatedAt.Time,
		})
	}

	// enrich (likes, liked, reply_count)
	enriched, err := s.EnrichPosts(ctx, viewerID, posts)
	if err != nil {
		return nil, err
	}

	// next cursor
	var nextCursor *time.Time
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		nextCursor = &last.CreatedAt
	}

	hasMore := len(posts) == int(limit)

	return &PaginatedPosts{
		Data:       enriched,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *Service) GetUserReplies(
	ctx context.Context,
	userID string,
	cursor *time.Time,
	limit int32,
	viewerID string,
) (*PaginatedPosts, error) {
	rows, err := s.repo.GetRepliesByUserWithCursor(ctx, userID, cursor, limit)
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, 0, len(rows))
	for _, r := range rows {
		posts = append(posts, &Post{
			ID:        r.ID,
			UserID:    r.UserID,
			Content:   r.Content,
			CreatedAt: r.CreatedAt.Time,
		})
	}

	enriched, err := s.EnrichPosts(ctx, viewerID, posts)
	if err != nil {
		return nil, err
	}

	var nextCursor *time.Time
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		nextCursor = &last.CreatedAt
	}

	hasMore := len(posts) == int(limit)

	return &PaginatedPosts{
		Data:       enriched,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
