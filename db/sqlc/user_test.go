package db

import (
	"context"
	"testing"

	_ "github.com/lib/pq"

	"github.com/mimzeslami/ExpenseShare/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Email:        util.RandomEmail(),
		PasswordHash: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.NotZero(t, user.ID)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
}
