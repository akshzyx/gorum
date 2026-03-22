-- name: CreateVerificationToken :exec
INSERT INTO email_verification_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: GetVerificationToken :one
SELECT token, user_id, expires_at, used, created_at
FROM email_verification_tokens
WHERE token = $1;

-- name: MarkVerificationTokenUsed :exec
UPDATE email_verification_tokens
SET used = TRUE
WHERE token = $1;
