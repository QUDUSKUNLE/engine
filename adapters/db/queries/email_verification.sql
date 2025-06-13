-- name: CreateEmailVerificationToken :one
INSERT INTO email_verification_tokens (
    email,
    token,
    expires_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetEmailVerificationToken :one
SELECT * FROM email_verification_tokens
WHERE email = $1 AND token = $2 AND used = false AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkEmailVerificationTokenUsed :exec
UPDATE email_verification_tokens
SET used = true
WHERE email = $1 AND token = $2;

-- name: MarkEmailAsVerified :exec
UPDATE users
SET email_verified = true,
    email_verified_at = NOW()
WHERE email = $1;
