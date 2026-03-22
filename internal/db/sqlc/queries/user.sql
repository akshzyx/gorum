-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash)
VALUES ($1, $2, $3, $4);

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, is_verified, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, username, email, password_hash, is_verified, created_at
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, is_verified, created_at
FROM users
WHERE username = $1;

-- name: ActivateUser :exec
UPDATE users
SET is_verified = TRUE
WHERE id = $1;

-- name: GetPublicProfileByUsername :one
SELECT id, username, created_at
FROM users
WHERE username = $1;

-- -- name: UpdateUserProfile :exec
-- UPDATE users
-- SET bio = $2,
--     avatar_url = $3
-- WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2,
    is_verified = FALSE
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;
