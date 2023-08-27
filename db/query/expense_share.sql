-- expense_shares.sql

-- Create an expense share
-- name: CreateExpenseShare :one
INSERT INTO expense_shares (
  expense_id,
  user_id,
  share,
  paid_status
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- Get an expense share by ID
-- name: GetExpenseShareByID :one
SELECT * FROM expense_shares
WHERE id = $1 LIMIT 1;

-- List expense shares for an expense with pagination
-- name: ListExpenseShares :many
SELECT * FROM expense_shares
WHERE expense_id = $1
LIMIT $2 OFFSET $3;

-- Update an expense share by ID
-- name: UpdateExpenseShare :one
UPDATE expense_shares SET
  share = $2,
  paid_status = $3
WHERE id = $1 RETURNING *;

-- Delete an expense share by ID
-- name: DeleteExpenseShare :exec
DELETE FROM expense_shares
WHERE id = $1;
