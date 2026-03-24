-- name: CreateFollow :exec
INSERT INTO follows (follower_id, following_id)
VALUES ($1, $2);

-- name: DeleteFollow :exec
DELETE FROM follows
WHERE follower_id = $1 AND following_id = $2;

-- name: GetFollowers :many
SELECT u.*
FROM users u
JOIN follows f ON u.id = f.follower_id
WHERE f.following_id = $1
ORDER BY f.created_at DESC;

-- name: GetFollowing :many
SELECT u.*
FROM users u
JOIN follows f ON u.id = f.following_id
WHERE f.follower_id = $1
ORDER BY f.created_at DESC;