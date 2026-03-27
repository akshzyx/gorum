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
SELECT 
    p.id,
    p.user_id,
    p.content,
    p.parent_post_id,
    p.created_at,
    u.username,

    -- total likes
    COALESCE(l.likes, 0) AS likes,

    -- whether current user liked
    CASE 
        WHEN ul.post_id IS NOT NULL THEN true
        ELSE false
    END AS liked

FROM posts p
JOIN users u ON u.id = p.user_id

-- likes count
LEFT JOIN (
    SELECT post_id, COUNT(*) AS likes
    FROM post_likes
    GROUP BY post_id
) l ON l.post_id = p.id

-- user liked
LEFT JOIN post_likes ul 
    ON ul.post_id = p.id 
    AND ul.user_id = sqlc.arg(user_id)

WHERE (p.id = $1 OR p.root_post_id = $1)
AND p.deleted_at IS NULL

ORDER BY p.created_at ASC;