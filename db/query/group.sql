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
WHERE id = $1 AND created_by_id =$2 LIMIT 1;

-- List groups with pagination
-- name: ListGroups :many
SELECT * FROM groups
WHERE created_by_id = $1
LIMIT $2 OFFSET $3;

-- Update a group by ID
-- name: UpdateGroup :one
UPDATE groups SET
  name = $1,
  category_id = $2,
  image_path = $3
WHERE id = $4 AND created_by_id =$5 RETURNING *;

-- Delete a group by ID
-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1 AND created_by_id =$2;
