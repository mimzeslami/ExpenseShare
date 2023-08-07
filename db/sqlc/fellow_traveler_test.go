package db

import (
	"context"
	"testing"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomFellowTraveller(t *testing.T) FellowTravelers {
	trip := createRandomTrip(t)
	arg := CreateFellowTravelersParams{
		TripID:          trip.ID,
		FellowFirstName: util.RandomString(6),
		FellowLastName:  util.RandomString(6),
	}
	fellowTraveller, err := testQueries.CreateFellowTravelers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, fellowTraveller)
	require.Equal(t, arg.TripID, fellowTraveller.TripID)
	require.Equal(t, arg.FellowFirstName, fellowTraveller.FellowFirstName)
	require.Equal(t, arg.FellowLastName, fellowTraveller.FellowLastName)
	require.NotZero(t, fellowTraveller.ID)

	return fellowTraveller

}

func TestCreateFellowTravelers(t *testing.T) {
	createRandomFellowTraveller(t)
}

func TestGetFellowTravelers(t *testing.T) {
	fellowTraveller1 := createRandomFellowTraveller(t)
	fellowTraveller2, err := testQueries.GetFellowTraveler(context.Background(), fellowTraveller1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fellowTraveller2)
	require.Equal(t, fellowTraveller1.ID, fellowTraveller2.ID)
	require.Equal(t, fellowTraveller1.TripID, fellowTraveller2.TripID)
	require.Equal(t, fellowTraveller1.FellowFirstName, fellowTraveller2.FellowFirstName)
	require.Equal(t, fellowTraveller1.FellowLastName, fellowTraveller2.FellowLastName)
}

func TestGetTripFellowTravelers(t *testing.T) {
	fellowTraveller1 := createRandomFellowTraveller(t)
	fellowTraveller2, err := testQueries.GetTripFellowTravelers(context.Background(), fellowTraveller1.TripID)
	require.NoError(t, err)
	require.NotEmpty(t, fellowTraveller2)
	require.Equal(t, fellowTraveller1.ID, fellowTraveller2[0].ID)
	require.Equal(t, fellowTraveller1.TripID, fellowTraveller2[0].TripID)
	require.Equal(t, fellowTraveller1.FellowFirstName, fellowTraveller2[0].FellowFirstName)
	require.Equal(t, fellowTraveller1.FellowLastName, fellowTraveller2[0].FellowLastName)
}
