-- Groups

-- Create a group
-- name: CreateGroup :one
INSERT INTO groups (
  name,
  category_id,
  created_by_id,
  image_path
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- Get a group by ID
-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id = $1 LIMIT 1;

-- List groups with pagination
-- name: ListGroups :many
SELECT * FROM groups
LIMIT $1 OFFSET $2;

-- Update a group by ID
-- name: UpdateGroup :one
UPDATE groups SET
  name = $2,
  category_id = $3,
  created_by_id = $4
WHERE id = $1 RETURNING *;

-- Delete a group by ID
-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;
