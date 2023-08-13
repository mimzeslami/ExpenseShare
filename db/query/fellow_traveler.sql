-- name: CreateFellowTravelers :one
INSERT INTO fellow_travelers (
  trip_id,
  fellow_first_name,
  fellow_last_name
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetFellowTraveler :one
SELECT fellow_travelers.* FROM fellow_travelers
LEFT JOIN trips ON fellow_travelers.trip_id = trips.id
WHERE fellow_travelers.id = $1 AND trips.user_id = $2 LIMIT 1;



-- name: GetTripFellowTravelers :many
SELECT fellow_travelers.*  FROM fellow_travelers
LEFT JOIN trips ON fellow_travelers.trip_id = trips.id
WHERE fellow_travelers.trip_id = $1 AND trips.user_id = $2;


-- name: UpdateFellowTraveler :one
UPDATE fellow_travelers SET
  fellow_first_name = $2,
  fellow_last_name = $3
WHERE id = $1 RETURNING *;

-- name: DeleteFellowTraveler :exec
DELETE FROM fellow_travelers WHERE id = $1;


-- name: DeleteTripFellowTravelers :exec
DELETE FROM fellow_travelers
WHERE trip_id = $1;



