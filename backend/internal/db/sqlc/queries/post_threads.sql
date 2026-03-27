-- name: GetPostForReply :one
SELECT id, root_post_id
FROM posts
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateReply :one
INSERT INTO posts (id, user_id, content, parent_post_id, root_post_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, content, created_at;

-- name: ListReplies :many
SELECT 
    posts.id,
    posts.user_id,
    posts.content,
    posts.created_at,
    users.username
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.parent_post_id = $1 
AND posts.deleted_at IS NULL
ORDER BY posts.created_at ASC;

-- name: GetThread :many
SELECT id, user_id, content, parent_post_id, created_at
FROM posts
WHERE (id = $1 OR root_post_id = $1)
AND deleted_at IS NULL
ORDER BY created_at ASC;
