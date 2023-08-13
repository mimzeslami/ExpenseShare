-- name: CreateExpense :one
INSERT INTO expenses (
  trip_id,
  payer_traveler_id,
  amount,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetExpense :one
SELECT expenses.* FROM expenses
LEFT JOIN trips ON expenses.trip_id = trips.id
WHERE expenses.id = $1 AND trips.user_id = $2 LIMIT 1;


-- name: GetTripExpenses :many
SELECT expenses.*  FROM expenses
LEFT JOIN trips ON expenses.trip_id = trips.id
WHERE expenses.trip_id = $1 AND trips.user_id = $2;


-- name: UpdateExpense :one
UPDATE expenses SET
  trip_id = $1,
  payer_traveler_id = $2,
  amount = $3,
  description = $4
WHERE id = $5 RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1;


-- name: DeleteTripExpenses :exec
DELETE FROM expenses
WHERE trip_id = $1;








