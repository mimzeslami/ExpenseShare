package db

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomNotification(t *testing.T, user Users) Notifications {
	arg := CreateNotificationParams{
		UserID:  user.ID,
		Message: util.RandomString(6),
		IsRead:  util.RandomBool(),
	}
	notification, err := testQueries.CreateNotification(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, notification)

	require.Equal(t, arg.UserID, notification.UserID)
	require.Equal(t, arg.Message, notification.Message)
	require.NotZero(t, notification.ID)
	require.NotZero(t, notification.CreatedAt)

	return notification
}

func TestCreateNotification(t *testing.T) {
	user := createRandomUser(t)
	createRandomNotification(t, user)
}

func TestDeleteNotification(t *testing.T) {
	user := createRandomUser(t)
	notification := createRandomNotification(t, user)
	err := testQueries.DeleteNotification(context.Background(), notification.ID)
	require.NoError(t, err)
}

func TestGetNotificationByID(t *testing.T) {
	user := createRandomUser(t)
	notification1 := createRandomNotification(t, user)
	notification2, err := testQueries.GetNotificationByID(context.Background(), notification1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, notification2)

	require.Equal(t, notification1.ID, notification2.ID)
	require.Equal(t, notification1.UserID, notification2.UserID)
	require.Equal(t, notification1.Message, notification2.Message)
	require.Equal(t, notification1.IsRead, notification2.IsRead)
	require.WithinDuration(t, notification1.CreatedAt, notification2.CreatedAt, time.Second)
}

func TestListNotifications(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomNotification(t, user)
	}

	arg := ListNotificationsParams{
		Limit:  5,
		Offset: 5,
		UserID: user.ID,
	}

	notifications, err := testQueries.ListNotifications(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, notifications, 5)

	for _, notification := range notifications {
		require.NotEmpty(t, notification)
	}
}

func TestMarkNotificationAsRead(t *testing.T) {
	user := createRandomUser(t)
	notification := createRandomNotification(t, user)
	err := testQueries.MarkNotificationAsRead(context.Background(), notification.ID)
	require.NoError(t, err)

	notification2, err := testQueries.GetNotificationByID(context.Background(), notification.ID)
	require.NoError(t, err)
	require.NotEmpty(t, notification2)

	require.Equal(t, notification.ID, notification2.ID)
	require.Equal(t, notification.UserID, notification2.UserID)
	require.Equal(t, notification.Message, notification2.Message)
	require.Equal(t, true, notification2.IsRead)
	require.WithinDuration(t, notification.CreatedAt, notification2.CreatedAt, time.Second)

}
