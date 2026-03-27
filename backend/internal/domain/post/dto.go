package post

type CreatePostRequest struct {
	Content string `json:"content" validate:"required,min=1,max=420"`
}

type CreatePostResponse struct {
	ID string `json:"id"`
}

type PublicPostResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CreateReplyRequest struct {
	Content string `json:"content" validate:"required,min=1,max=420"`
}

type ReplyResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	Username     string  `json:"username"`
	Content      string  `json:"content"`
	CreatedAt    string  `json:"created_at"`
	ParentPostID *string `json:"parent_post_id,omitempty"`
}
