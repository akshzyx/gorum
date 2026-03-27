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

	// fetch root info also
	parent, err := s.repo.GetPostForReply(ctx, id)
	if err == nil && parent.RootPostID != nil {
		post.RootPostID = parent.RootPostID
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

func (s *Service) ListLatest(ctx context.Context, limit int32, cursor string) ([]*Post, string, error) {
	var cursorTime *time.Time
	var cursorID *string

	if cursor != "" {
		parts := strings.Split(cursor, "|")

		if len(parts) == 2 {
			t, err := time.Parse(time.RFC3339, parts[0])
			if err == nil {
				cursorTime = &t
				cursorID = &parts[1]
			}
		}
	}

	posts, err := s.repo.ListLatest(ctx, limit+1, cursorTime, cursorID)
	if err != nil {
		return nil, "", err
	}

	hasMore := false
	if int32(len(posts)) > limit {
		hasMore = true
		posts = posts[:limit]
	}

	var nextCursor string
	if hasMore && len(posts) > 0 {
		last := posts[len(posts)-1]
		nextCursor = last.CreatedAt.Format(time.RFC3339) + "|" + last.ID
	}

	return posts, nextCursor, nil
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

	// prevent duplicate likes (safety at service layer)
	liked, err := s.repo.HasUserLiked(ctx, userID, postID)
	if err == nil && liked {
		return nil
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

	// batch replies count (optimized)
	replyCountMap, err := s.repo.GetRepliesCountByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	resp := make([]map[string]interface{}, 0, len(posts))

	for _, p := range posts {
		resp = append(resp, map[string]interface{}{
			"id":          p.ID,
			"user_id":     p.UserID,
			"username":    p.Username,
			"content":     p.Content,
			"created_at":  p.CreatedAt,
			"likes":       likesMap[p.ID],
			"liked":       likedMap[p.ID],
			"reply_count": replyCountMap[p.ID],
		})
	}

	return resp, nil
}

func (s *Service) GetUserPosts(ctx context.Context, userID string, limit int32) ([]*Post, error) {
	return s.repo.GetPostsByUser(ctx, userID, limit)
}

func (s *Service) GetUserReplies(ctx context.Context, userID string, limit int32) ([]*Post, error) {
	return s.repo.GetRepliesByUser(ctx, userID, limit)
}

func (s *Service) GetUserPostsPaginated(ctx context.Context, userID string, limit int32, cursor string) ([]*Post, string, error) {
	var cursorTime *time.Time

	if cursor != "" {
		t, err := time.Parse(time.RFC3339, cursor)
		if err != nil {
			return nil, "", err
		}
		cursorTime = &t
	}

	posts, err := s.repo.GetPostsByUserCursor(ctx, userID, limit+1, cursorTime)
	if err != nil {
		return nil, "", err
	}

	hasMore := false
	if int32(len(posts)) > limit {
		hasMore = true
		posts = posts[:limit]
	}

	var nextCursor string
	if hasMore && len(posts) > 0 {
		nextCursor = posts[len(posts)-1].CreatedAt.Format(time.RFC3339)
	}

	return posts, nextCursor, nil
}

func (s *Service) GetUserRepliesPaginated(ctx context.Context, userID string, limit int32, cursor string) ([]*Post, string, error) {
	var cursorTime *time.Time

	if cursor != "" {
		t, err := time.Parse(time.RFC3339, cursor)
		if err != nil {
			return nil, "", err
		}
		cursorTime = &t
	}

	posts, err := s.repo.GetRepliesByUserCursor(ctx, userID, limit+1, cursorTime)
	if err != nil {
		return nil, "", err
	}

	hasMore := false
	if int32(len(posts)) > limit {
		hasMore = true
		posts = posts[:limit]
	}

	var nextCursor string
	if hasMore && len(posts) > 0 {
		nextCursor = posts[len(posts)-1].CreatedAt.Format(time.RFC3339)
	}

	return posts, nextCursor, nil
}

func (s *Service) ListRepliesPaginated(ctx context.Context, postID string, limit int32, cursor string, order string) ([]*Post, string, error) {
	var cursorTime *time.Time

	if cursor != "" {
		t, err := time.Parse(time.RFC3339, cursor)
		if err != nil {
			return nil, "", err
		}
		cursorTime = &t
	}

	var posts []*Post
	var err error

	if order == "desc" {
		posts, err = s.repo.ListRepliesCursorDesc(ctx, postID, limit+1, cursorTime)
	} else {
		posts, err = s.repo.ListRepliesCursorAsc(ctx, postID, limit+1, cursorTime)
	}

	if err != nil {
		return nil, "", err
	}

	hasMore := false
	if int32(len(posts)) > limit {
		hasMore = true
		posts = posts[:limit]
	}

	var nextCursor string
	if hasMore && len(posts) > 0 {
		nextCursor = posts[len(posts)-1].CreatedAt.Format(time.RFC3339)
	}

	return posts, nextCursor, nil
}
