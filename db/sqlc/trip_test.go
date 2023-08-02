package db

import (
	"context"
	"testing"

	"github.com/mimzeslami/ExpenseShare/util"
	"github.com/stretchr/testify/require"
)

func createRandomTrip(t *testing.T) Trips {
	user := createRandomUser(t)
	arg := CreateTripParams{
		TripName:  util.RandomString(6),
		StartDate: util.RandomDatetime(),
		EndDate:   util.RandomDatetime(),
		UserID:    user.ID,
	}
	trip, err := testQueries.CreateTrip(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trip)
	require.Equal(t, arg.TripName, trip.TripName)
	require.Equal(t, arg.UserID, trip.UserID)
	require.NotZero(t, trip.ID)
	return trip

}

func TestCreateTrip(t *testing.T) {
	createRandomTrip(t)
}