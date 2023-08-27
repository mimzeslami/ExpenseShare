package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/mimzeslami/expense_share/util"
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
		Phone:        util.RandomString(9),
		ImagePath:    util.RandomString(6),
		TimeZone:     util.RandomString(6),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.ImagePath, user.ImagePath)
	require.Equal(t, arg.TimeZone, user.TimeZone)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserById(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.ImagePath, user2.ImagePath)
	require.Equal(t, user1.TimeZone, user2.TimeZone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.ImagePath, user2.ImagePath)
	require.Equal(t, user1.TimeZone, user2.TimeZone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)

	}
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		ID:           user1.ID,
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Email:        util.RandomEmail(),
		PasswordHash: util.RandomString(6),
		Phone:        util.RandomString(9),
		ImagePath:    util.RandomString(6),
		TimeZone:     util.RandomString(6),
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, arg.ID, user2.ID)
	require.Equal(t, arg.FirstName, user2.FirstName)
	require.Equal(t, arg.LastName, user2.LastName)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, arg.PasswordHash, user2.PasswordHash)
	require.Equal(t, arg.Phone, user2.Phone)
	require.Equal(t, arg.ImagePath, user2.ImagePath)
	require.Equal(t, arg.TimeZone, user2.TimeZone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
