-- name: CreateTrip :one
INSERT INTO trips (
  title,
  start_date,
  end_date,
  user_id
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTrip :one
SELECT * FROM trips
WHERE id = $1 LIMIT 1;

-- name: ListTrip :many
SELECT * FROM trips
WHERE user_id = $1
LIMIT $2 OFFSET $3;

-- name: UpdateTrip :one
UPDATE trips SET
  title = $2,
  start_date = $3,
  end_date = $4
WHERE id = $1 RETURNING *;

-- name: DeleteTrip :exec
DELETE FROM trips
WHERE id = $1;






