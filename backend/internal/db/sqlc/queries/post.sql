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
SELECT 
  p.id,
  p.user_id,
  p.content,
  p.created_at,
  u.username
FROM posts p
JOIN users u ON u.id = p.user_id
WHERE p.deleted_at IS NULL
AND p.parent_post_id IS NULL
AND ($1::timestamptz IS NULL OR p.created_at < $1)
ORDER BY p.created_at DESC
LIMIT $2;

-- name: CountReplies :one
SELECT COUNT(*)
FROM posts
WHERE root_post_id = $1
AND deleted_at IS NULL;

-- name: CreateLike :exec
INSERT INTO post_likes (user_id, post_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteLike :exec
DELETE FROM post_likes
WHERE user_id = $1 AND post_id = $2;

-- name: CountLikes :one
SELECT COUNT(*) FROM post_likes
WHERE post_id = $1;

-- name: HasUserLikedPost :one
SELECT EXISTS (
    SELECT 1 FROM post_likes
    WHERE user_id = $1 AND post_id = $2
);

-- name: GetLikesCountByPostIDs :many
SELECT post_id, COUNT(*) AS count
FROM post_likes
WHERE post_id = ANY($1::text[])
GROUP BY post_id;

-- name: GetUserLikedPosts :many
SELECT post_id
FROM post_likes
WHERE user_id = $1
AND post_id = ANY(sqlc.arg(post_ids)::text[]);

-- name: GetPostsByUser :many
SELECT id, user_id, content, created_at
FROM posts
WHERE user_id = $1
AND parent_post_id IS NULL
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2;

-- name: GetRepliesByUser :many
SELECT id, user_id, content, created_at
FROM posts
WHERE user_id = $1
AND parent_post_id IS NOT NULL
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2;

-- name: GetRepliesCountByPostIDs :many
SELECT parent_post_id AS post_id, COUNT(*) AS count
FROM posts
WHERE parent_post_id = ANY($1::text[])
GROUP BY parent_post_id;

-- name: GetPostsByUserCursor :many
SELECT id, user_id, content, created_at
FROM posts
WHERE user_id = $1
AND parent_post_id IS NULL
AND deleted_at IS NULL
AND ($2::timestamptz IS NULL OR created_at < $2)
ORDER BY created_at DESC
LIMIT $3;

-- name: GetRepliesByUserCursor :many
SELECT id, user_id, content, created_at
FROM posts
WHERE user_id = $1
AND parent_post_id IS NOT NULL
AND deleted_at IS NULL
AND ($2::timestamptz IS NULL OR created_at < $2)
ORDER BY created_at DESC
LIMIT $3;

-- name: ListRepliesCursorAsc :many
SELECT id, user_id, content, created_at
FROM posts
WHERE parent_post_id = $1
AND deleted_at IS NULL
AND ($2::timestamptz IS NULL OR created_at > $2)
ORDER BY created_at ASC
LIMIT $3;

-- name: ListRepliesCursorDesc :many
SELECT id, user_id, content, created_at
FROM posts
WHERE parent_post_id = $1
AND deleted_at IS NULL
AND ($2::timestamptz IS NULL OR created_at < $2)
ORDER BY created_at DESC
LIMIT $3;