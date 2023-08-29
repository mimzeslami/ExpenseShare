-- expenses.sql

-- Create an expense
-- name: CreateExpense :one
INSERT INTO expenses (
  group_id,
  paid_by_id,
  amount,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- Get an expense by ID
-- name: GetExpenseByID :one
SELECT * FROM expenses
WHERE id = $1 LIMIT 1;

-- List expenses for a group with pagination
-- name: ListExpenses :many
SELECT * FROM expenses
WHERE group_id = $1
LIMIT $2 OFFSET $3;

-- Update an expense by ID
-- name: UpdateExpense :one
UPDATE expenses SET
  amount = $2,
  description = $3
WHERE id = $1 RETURNING *;

-- Delete an expense by ID
-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1;
