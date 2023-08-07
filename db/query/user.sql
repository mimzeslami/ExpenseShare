-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  password_hash,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM users
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users SET
  first_name = $2,
  last_name = $3,
  password_hash = $4,
  email = $5
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;




