// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: fellow_traveler.sql

package db

import (
	"context"
)

const createFellowTravelers = `-- name: CreateFellowTravelers :one
INSERT INTO fellow_travelers (
  trip_id,
  fellow_first_name,
  fellow_last_name
) VALUES (
  $1, $2, $3
) RETURNING id, trip_id, fellow_first_name, fellow_last_name, created_at
`

type CreateFellowTravelersParams struct {
	TripID          int64  `json:"trip_id"`
	FellowFirstName string `json:"fellow_first_name"`
	FellowLastName  string `json:"fellow_last_name"`
}

func (q *Queries) CreateFellowTravelers(ctx context.Context, arg CreateFellowTravelersParams) (FellowTravelers, error) {
	row := q.db.QueryRowContext(ctx, createFellowTravelers, arg.TripID, arg.FellowFirstName, arg.FellowLastName)
	var i FellowTravelers
	err := row.Scan(
		&i.ID,
		&i.TripID,
		&i.FellowFirstName,
		&i.FellowLastName,
		&i.CreatedAt,
	)
	return i, err
}

const getFellowTraveler = `-- name: GetFellowTraveler :one
SELECT id, trip_id, fellow_first_name, fellow_last_name, created_at FROM fellow_travelers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFellowTraveler(ctx context.Context, id int64) (FellowTravelers, error) {
	row := q.db.QueryRowContext(ctx, getFellowTraveler, id)
	var i FellowTravelers
	err := row.Scan(
		&i.ID,
		&i.TripID,
		&i.FellowFirstName,
		&i.FellowLastName,
		&i.CreatedAt,
	)
	return i, err
}

const getTripFellowTravelers = `-- name: GetTripFellowTravelers :many
SELECT id, trip_id, fellow_first_name, fellow_last_name, created_at FROM fellow_travelers
WHERE trip_id = $1
`

func (q *Queries) GetTripFellowTravelers(ctx context.Context, tripID int64) ([]FellowTravelers, error) {
	rows, err := q.db.QueryContext(ctx, getTripFellowTravelers, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FellowTravelers{}
	for rows.Next() {
		var i FellowTravelers
		if err := rows.Scan(
			&i.ID,
			&i.TripID,
			&i.FellowFirstName,
			&i.FellowLastName,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
