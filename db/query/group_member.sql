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


-- Delete all group members for a group
-- name: DeleteGroupMembers :exec
DELETE FROM group_members
WHERE group_id = $1;

-- List group members for a group with pagination, including user and group details
-- name: ListGroupMembersWithDetails :many
SELECT 
    gm.*,
    u.id AS user_id,
    u.first_name AS user_first_name,
    u.last_name AS user_last_name,
    u.email AS user_email,
    u.phone AS user_phone,
    g.id AS group_id,
    g.name AS group_name,
    g.category_id AS group_category_id,
    g.created_by_id AS group_created_by_id
FROM 
    group_members gm
JOIN 
    users u ON gm.user_id = u.id
JOIN 
    groups g ON gm.group_id = g.id
WHERE 
    gm.group_id = $1
LIMIT 
    $2 OFFSET $3;


-- Get Group Member by Group ID and User ID
-- name: GetGroupMemberByGroupIDAndUserID :one
SELECT * FROM group_members
WHERE group_id = $1 AND user_id = $2 LIMIT 1;

