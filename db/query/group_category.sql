-- name: CreateGroupCategory :one
INSERT INTO group_categories (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetGroupCategory :one
SELECT * FROM group_categories
WHERE id = $1 LIMIT 1;

-- name: ListGroupCategories :many
SELECT * FROM group_categories
LIMIT $1 OFFSET $2;

-- name: UpdateGroupCategory :one
UPDATE group_categories SET
  name = $2
WHERE id = $1 RETURNING *;

-- name: DeleteGroupCategory :exec
DELETE FROM group_categories
WHERE id = $1;
