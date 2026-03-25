package postgres

import (
	"context"
	"time"

	db "github.com/akshzyx/gorum/internal/db/sqlc/generated"
	"github.com/akshzyx/gorum/internal/domain/post"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostRepository struct {
	q *db.Queries
}

func NewPostRepository(q *db.Queries) *PostRepository {
	return &PostRepository{q: q}
}

func (r *PostRepository) Create(ctx context.Context, p *post.Post) error {
	_, err := r.q.CreatePost(ctx, db.CreatePostParams{
		ID:      p.ID,
		UserID:  p.UserID,
		Content: p.Content,
	})
	return err
}

func (r *PostRepository) GetByID(ctx context.Context, id string) (*post.Post, error) {
	row, err := r.q.GetPostByID(ctx, id)
	if err != nil {
		return nil, post.ErrPostNotFound
	}

	return &post.Post{
		ID:        row.ID,
		UserID:    row.UserID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

// Delete by post ID. Ownership is checked in service layer.
func (r *PostRepository) Delete(ctx context.Context, postID string) error {
	return r.q.DeletePostByID(ctx, postID)
}

func (r *PostRepository) ListLatest(ctx context.Context, limit int32) ([]*post.Post, error) {
	rows, err := r.q.ListLatestPosts(ctx, limit)
	if err != nil {
		return nil, err
	}

	posts := make([]*post.Post, 0, len(rows))
	for _, row := range rows {
		posts = append(posts, &post.Post{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return posts, nil
}

func (r *PostRepository) GetPostForReply(ctx context.Context, id string) (*post.Post, error) {
	row, err := r.q.GetPostForReply(ctx, id)
	if err != nil {
		return nil, post.ErrPostNotFound
	}

	var rootID *string
	if row.RootPostID.Valid {
		rootID = &row.RootPostID.String
	}

	return &post.Post{
		ID:         row.ID,
		RootPostID: rootID,
	}, nil
}

func (r *PostRepository) CreateReply(ctx context.Context, p *post.Post) error {
	_, err := r.q.CreateReply(ctx, db.CreateReplyParams{
		ID:           p.ID,
		UserID:       p.UserID,
		Content:      p.Content,
		ParentPostID: pgtype.Text{String: *p.ParentPostID, Valid: true},
		RootPostID:   pgtype.Text{String: *p.RootPostID, Valid: true},
	})
	return err
}

func (r *PostRepository) ListReplies(ctx context.Context, postID string) ([]*post.Post, error) {
	rows, err := r.q.ListReplies(ctx, pgtype.Text{
		String: postID,
		Valid:  true,
	})
	if err != nil {
		return nil, err
	}

	var posts []*post.Post
	for _, row := range rows {
		posts = append(posts, &post.Post{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return posts, nil
}

func (r *PostRepository) CountReplies(ctx context.Context, postID string) (int64, error) {
	return r.q.CountReplies(ctx, pgtype.Text{
		String: postID,
		Valid:  true,
	})
}

func (r *PostRepository) GetThread(ctx context.Context, rootID string) ([]*post.Post, error) {
	rows, err := r.q.GetThread(ctx, rootID)
	if err != nil {
		return nil, err
	}

	var posts []*post.Post
	for _, row := range rows {
		var parentID *string
		if row.ParentPostID.Valid {
			parentID = &row.ParentPostID.String
		}

		posts = append(posts, &post.Post{
			ID:           row.ID,
			UserID:       row.UserID,
			Content:      row.Content,
			ParentPostID: parentID,
			CreatedAt:    row.CreatedAt.Time,
		})
	}

	return posts, nil
}

// Like system
func (r *PostRepository) CreateLike(ctx context.Context, userID, postID string) error {
	return r.q.CreateLike(ctx, db.CreateLikeParams{
		UserID: userID,
		PostID: postID,
	})
}

func (r *PostRepository) DeleteLike(ctx context.Context, userID, postID string) error {
	return r.q.DeleteLike(ctx, db.DeleteLikeParams{
		UserID: userID,
		PostID: postID,
	})
}

func (r *PostRepository) CountLikes(ctx context.Context, postID string) (int64, error) {
	return r.q.CountLikes(ctx, postID)
}

func (r *PostRepository) HasUserLiked(ctx context.Context, userID, postID string) (bool, error) {
	return r.q.HasUserLikedPost(ctx, db.HasUserLikedPostParams{
		UserID: userID,
		PostID: postID,
	})
}

func (r *PostRepository) GetLikesCountByPostIDs(ctx context.Context, postIDs []string) (map[string]int64, error) {
	rows, err := r.q.GetLikesCountByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)
	for _, row := range rows {
		result[row.PostID] = row.Count
	}

	return result, nil
}

func (r *PostRepository) GetUserLikedPosts(ctx context.Context, userID string, postIDs []string) (map[string]bool, error) {
	rows, err := r.q.GetUserLikedPosts(ctx, db.GetUserLikedPostsParams{
		UserID:  userID,
		PostIds: postIDs,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, postID := range rows {
		result[postID] = true
	}

	return result, nil
}

func (r *PostRepository) GetPostsByUser(ctx context.Context, userID string, limit int32) ([]*post.Post, error) {
	rows, err := r.q.GetPostsByUser(ctx, db.GetPostsByUserParams{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	var posts []*post.Post
	for _, row := range rows {
		posts = append(posts, &post.Post{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return posts, nil
}

func (r *PostRepository) GetRepliesByUser(ctx context.Context, userID string, limit int32) ([]*post.Post, error) {
	rows, err := r.q.GetRepliesByUser(ctx, db.GetRepliesByUserParams{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	var posts []*post.Post
	for _, row := range rows {
		posts = append(posts, &post.Post{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return posts, nil
}

func (r *PostRepository) ListLatestWithCursor(
	ctx context.Context,
	cursor *time.Time,
	limit int32,
) ([]db.ListLatestPostsWithCursorRow, error) {
	var pgCursor pgtype.Timestamp

	if cursor != nil {
		pgCursor = pgtype.Timestamp{
			Time:  *cursor,
			Valid: true,
		}
	} else {
		pgCursor = pgtype.Timestamp{
			Valid: false,
		}
	}

	return r.q.ListLatestPostsWithCursor(ctx, db.ListLatestPostsWithCursorParams{
		Column1: pgCursor,
		Limit:   limit,
	})
}

func (r *PostRepository) ListRepliesWithCursor(
	ctx context.Context,
	rootID string,
	cursor *time.Time,
	limit int32,
) ([]db.ListRepliesWithCursorRow, error) {
	var pgCursor pgtype.Timestamp

	if cursor != nil {
		pgCursor = pgtype.Timestamp{
			Time:  *cursor,
			Valid: true,
		}
	} else {
		pgCursor = pgtype.Timestamp{
			Valid: false,
		}
	}

	return r.q.ListRepliesWithCursor(ctx, db.ListRepliesWithCursorParams{
		RootPostID: pgtype.Text{
			String: rootID,
			Valid:  true,
		},
		Column2: pgCursor,
		Limit:   limit,
	})
}

func (r *PostRepository) GetPostsByUserWithCursor(
	ctx context.Context,
	userID string,
	cursor *time.Time,
	limit int32,
) ([]db.GetPostsByUserWithCursorRow, error) {
	var pgCursor pgtype.Timestamp

	if cursor != nil {
		pgCursor = pgtype.Timestamp{
			Time:  *cursor,
			Valid: true,
		}
	} else {
		pgCursor = pgtype.Timestamp{
			Valid: false,
		}
	}

	return r.q.GetPostsByUserWithCursor(ctx, db.GetPostsByUserWithCursorParams{
		UserID:  userID,
		Column2: pgCursor,
		Limit:   limit,
	})
}

func (r *PostRepository) GetRepliesByUserWithCursor(
	ctx context.Context,
	userID string,
	cursor *time.Time,
	limit int32,
) ([]db.GetRepliesByUserWithCursorRow, error) {
	var pgCursor pgtype.Timestamp

	if cursor != nil {
		pgCursor = pgtype.Timestamp{
			Time:  *cursor,
			Valid: true,
		}
	} else {
		pgCursor = pgtype.Timestamp{
			Valid: false,
		}
	}

	return r.q.GetRepliesByUserWithCursor(ctx, db.GetRepliesByUserWithCursorParams{
		UserID:  userID,
		Column2: pgCursor,
		Limit:   limit,
	})
}
