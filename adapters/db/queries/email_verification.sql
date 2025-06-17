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
WHERE token = $1 AND used = false
LIMIT 1;

-- name: MarkEmailVerificationTokenUsed :exec
UPDATE email_verification_tokens
SET used = true
WHERE id = $1;


UPDATE users
SET email_verified = true,
    email_verified_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkEmailAsVerified :exec
UPDATE users
SET email_verified = true,
    email_verified_at = NOW(),
    updated_at = NOW()
WHERE email = $1
RETURNING *;
