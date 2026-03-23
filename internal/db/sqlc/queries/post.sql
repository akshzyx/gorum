-- name: CreatePost :one
INSERT INTO posts (id, user_id, content)
VALUES ($1, $2, $3)
RETURNING id, user_id, content, created_at;

-- name: GetPostByID :one
SELECT id, user_id, content, created_at
FROM posts
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeletePostByID :exec
UPDATE posts
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLatestPosts :many
SELECT id, user_id, content, created_at
FROM posts
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1;