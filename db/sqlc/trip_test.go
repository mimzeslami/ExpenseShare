package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomTrip(t *testing.T, user Users) Trips {

	arg := CreateTripParams{
		Title:     util.RandomString(6),
		StartDate: util.RandomDatetime(),
		EndDate:   util.RandomDatetime(),
		UserID:    user.ID,
	}
	trip, err := testQueries.CreateTrip(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trip)
	require.Equal(t, arg.Title, trip.Title)
	require.Equal(t, arg.UserID, trip.UserID)
	require.NotZero(t, trip.ID)
	return trip

}

func TestCreateTrip(t *testing.T) {
	user := createRandomUser(t)
	createRandomTrip(t, user)
}

func TestGetTrip(t *testing.T) {
	user := createRandomUser(t)
	trip1 := createRandomTrip(t, user)
	arg := GetTripParams{
		ID:     trip1.ID,
		UserID: user.ID,
	}
	trip2, err := testQueries.GetTrip(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trip2)
	require.Equal(t, trip1.Title, trip2.Title)
	require.Equal(t, trip1.UserID, trip2.UserID)
	require.Equal(t, trip1.ID, trip2.ID)
}

func TestListTrip(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomTrip(t, user)
	}
	arg := ListTripParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 5,
	}
	trips, err := testQueries.ListTrip(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, trips, 5)
	for _, trip := range trips {
		require.NotEmpty(t, trip)
	}
}

func TestUpdateTrip(t *testing.T) {
	user := createRandomUser(t)
	trip1 := createRandomTrip(t, user)
	arg := UpdateTripParams{
		ID:        trip1.ID,
		Title:     util.RandomString(6),
		StartDate: util.RandomDatetime(),
		EndDate:   util.RandomDatetime(),
		UserID:    user.ID,
	}
	trip2, err := testQueries.UpdateTrip(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trip2)
	require.Equal(t, trip1.ID, trip2.ID)
	require.Equal(t, arg.Title, trip2.Title)
	require.WithinDuration(t, arg.StartDate, trip2.StartDate, time.Hour)
	require.WithinDuration(t, arg.EndDate, trip2.EndDate, time.Hour)
}

func TestDeleteTrip(t *testing.T) {
	user := createRandomUser(t)
	trip := createRandomTrip(t, user)

	deleteArg := DeleteTripParams{
		ID:     trip.ID,
		UserID: user.ID,
	}

	err := testQueries.DeleteTrip(context.Background(), deleteArg)
	require.NoError(t, err)

	arg := GetTripParams{
		ID:     trip.ID,
		UserID: user.ID,
	}
	trip1, err := testQueries.GetTrip(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, trip1)

}
