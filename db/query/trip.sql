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
