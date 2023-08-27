-- users.sql

-- Create a user
-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  password_hash,
  phone,
  image_path,
  time_zone
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- Get a user by ID
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- Get a user by email
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- List users with pagination
-- name: ListUsers :many
SELECT * FROM users
LIMIT $1 OFFSET $2;

-- Update a user by ID
-- name: UpdateUser :one
UPDATE users SET
  first_name = $2,
  last_name = $3,
  email = $4,
  password_hash = $5,
  phone = $6,
  image_path = $7,
  time_zone = $8
WHERE id = $1 RETURNING *;

-- Delete a user by ID
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
