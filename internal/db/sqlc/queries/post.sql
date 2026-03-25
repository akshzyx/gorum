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
AND parent_post_id IS NULL
ORDER BY created_at DESC
LIMIT $1;

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

-- name: ListLatestPostsWithCursor :many
SELECT id, user_id, content, created_at
FROM posts
WHERE deleted_at IS NULL
AND parent_post_id IS NULL
AND ($1::timestamp IS NULL OR created_at < $1)
ORDER BY created_at DESC
LIMIT $2;

-- name: ListRepliesWithCursor :many
SELECT id, user_id, content, parent_post_id, root_post_id, created_at
FROM posts
WHERE root_post_id = $1
AND deleted_at IS NULL
AND ($2::timestamp IS NULL OR created_at > $2)
ORDER BY created_at ASC
LIMIT $3;

-- name: GetPostsByUserWithCursor :many
SELECT id, user_id, content, created_at
FROM posts
WHERE user_id = $1
AND parent_post_id IS NULL
AND deleted_at IS NULL
AND ($2::timestamp IS NULL OR created_at < $2)
ORDER BY created_at DESC
LIMIT $3;

-- name: GetRepliesByUserWithCursor :many
SELECT id, user_id, content, created_at
FROM posts
WHERE user_id = $1
AND parent_post_id IS NOT NULL
AND deleted_at IS NULL
AND ($2::timestamp IS NULL OR created_at < $2)
ORDER BY created_at DESC
LIMIT $3;