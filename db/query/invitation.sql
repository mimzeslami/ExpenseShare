-- invitations.sql

-- Create an invitation
-- name: CreateInvitation :one
INSERT INTO invitations (
  inviter_id,
  invitee_id,
  group_id,
  status,
  code
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- Get an invitation by ID
-- name: GetInvitationByID :one
SELECT * FROM invitations
WHERE id = $1 LIMIT 1;

-- List invitations for an invitee with pagination
-- name: ListInvitationsForInvitee :many
SELECT * FROM invitations
WHERE invitee_id = $1
LIMIT $2 OFFSET $3;

-- Update an invitation by ID
-- name: UpdateInvitation :one
UPDATE invitations SET
  status = $2,
  accepted_at = $3,
  rejected_at = $4
WHERE id = $1 RETURNING *;

-- Delete an invitation by ID
-- name: DeleteInvitation :exec
DELETE FROM invitations
WHERE id = $1;

-- Get Current Invitation By GroupID And InviteeID
-- name: GetCurrentInvitationByGroupIDAndInviteeID :one
SELECT * FROM invitations
WHERE group_id = $1 AND invitee_id = $2 AND inviter_id = $3;

-- Get User info by invitation code
-- name: GetUserInfoByInvitationCode :one
SELECT users.id, users.email, users.first_name, users.last_name, users.phone, users.time_zone FROM users
INNER JOIN invitations ON users.id = invitations.invitee_id
WHERE invitations.code = $1;


-- Get invitation by code
-- name: GetInvitationByCode :one
SELECT * FROM invitations
WHERE code = $1;
