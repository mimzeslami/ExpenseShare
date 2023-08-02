// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateFellowTravelers(ctx context.Context, arg CreateFellowTravelersParams) (FellowTravelers, error)
	CreateTrip(ctx context.Context, arg CreateTripParams) (Trips, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (Users, error)
	GetFellowTraveler(ctx context.Context, id uuid.UUID) (FellowTravelers, error)
	GetTrip(ctx context.Context, id uuid.UUID) (Trips, error)
	GetTripFellowTravelers(ctx context.Context, tripID uuid.UUID) ([]FellowTravelers, error)
	GetUser(ctx context.Context, email string) (Users, error)
}

var _ Querier = (*Queries)(nil)
