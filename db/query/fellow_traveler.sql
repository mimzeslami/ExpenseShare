-- name: CreateFellowTravelers :one
INSERT INTO fellow_travelers (
  trip_id,
  fellow_first_name,
  fellow_last_name
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetFellowTraveler :one
SELECT * FROM fellow_travelers
WHERE id = $1 LIMIT 1;


-- name: GetTripFellowTravelers :many
SELECT * FROM fellow_travelers
WHERE trip_id = $1;

