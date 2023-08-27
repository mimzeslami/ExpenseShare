-- GroupMembers

-- Create a group member
-- name: CreateGroupMember :one
INSERT INTO group_members (
  group_id,
  user_id
) VALUES (
  $1, $2
) RETURNING *;

-- Get a group member by ID
-- name: GetGroupMemberByID :one
SELECT * FROM group_members
WHERE id = $1 LIMIT 1;

-- List group members for a group with pagination
-- name: ListGroupMembers :many
SELECT * FROM group_members
WHERE group_id = $1
LIMIT $2 OFFSET $3;

-- Update a group member by ID
-- name: UpdateGroupMember :one
UPDATE group_members SET
  user_id = $2
WHERE id = $1 RETURNING *;

-- Delete a group member by ID
-- name: DeleteGroupMember :exec
DELETE FROM group_members
WHERE id = $1;
