-- name: GetUser :one
SELECT * FROM users where id = $1;

-- name: CreateUser :one
INSERT INTO users (
  email,
  nin,
  password,
  user_type,
  phone_number,
  email_verified,
  fullname,
  created_admin
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8::uuid
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  nin = COALESCE($2, nin),
  fullname = COALESCE($3, fullname),
  phone_number = COALESCE($4, phone_number),
  updated_at = NOW()
WHERE id = $1
RETURNING id, email, nin, user_type, fullname, phone_number, email_verified, email_verified_at, created_at, updated_at;
