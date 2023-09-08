-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password,
  full_name,
  username
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;