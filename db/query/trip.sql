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
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListTrip :many
SELECT * FROM trips
WHERE user_id = $1
LIMIT $2 OFFSET $3;

-- name: UpdateTrip :one
UPDATE trips SET
  title = $1,
  start_date = $2,
  end_date = $3
WHERE id = $4 AND user_id =$5 RETURNING *;

-- name: DeleteTrip :exec
DELETE FROM trips
WHERE id = $1 AND user_id = $2;






