-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (
  email, token, expires_at
) VALUES ($1, $2, $3);

-- name: GetPasswordResetToken :one
SELECT * FROM password_reset_tokens 
WHERE token = $1 AND used = false 
ORDER BY created_at DESC 
LIMIT 1;

-- name: MarkResetTokenUsed :exec
UPDATE password_reset_tokens
SET used = true
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE email = $1;
