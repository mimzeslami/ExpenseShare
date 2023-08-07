package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomFellowTraveller(t *testing.T) FellowTravelers {
	user := createRandomUser(t)
	trip := createRandomTrip(t, user)
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

func TestDeleteFeeellowTraveler(t *testing.T) {
	fellowTraveller1 := createRandomFellowTraveller(t)
	err := testQueries.DeleteFellowTraveler(context.Background(), fellowTraveller1.ID)
	require.NoError(t, err)
	fellowTraveller2, err := testQueries.GetFellowTraveler(context.Background(), fellowTraveller1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, fellowTraveller2)

}

func TestDeleteTripFellowTravelers(t *testing.T) {
	fellowTraveller1 := createRandomFellowTraveller(t)
	err := testQueries.DeleteTripFellowTravelers(context.Background(), fellowTraveller1.TripID)
	require.NoError(t, err)
	fellowTraveller2, err := testQueries.GetFellowTraveler(context.Background(), fellowTraveller1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, fellowTraveller2)
}

func TestUpdateFellowTraveler(t *testing.T) {
	fellowTraveller1 := createRandomFellowTraveller(t)
	arg := UpdateFellowTravelerParams{
		ID:              fellowTraveller1.ID,
		FellowFirstName: util.RandomString(6),
		FellowLastName:  util.RandomString(6),
	}
	fellowTraveller2, err := testQueries.UpdateFellowTraveler(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, fellowTraveller2)
	require.Equal(t, fellowTraveller1.ID, fellowTraveller2.ID)
	require.Equal(t, arg.FellowFirstName, fellowTraveller2.FellowFirstName)
	require.Equal(t, arg.FellowLastName, fellowTraveller2.FellowLastName)
}
