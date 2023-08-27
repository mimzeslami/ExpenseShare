package db

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomExpenseShare(t *testing.T, expense Expenses) ExpenseShares {
	user := createRandomUser(t)
	arg := CreateExpenseShareParams{
		ExpenseID:  expense.ID,
		UserID:     user.ID,
		Share:      util.RandomMoney(),
		PaidStatus: util.RandomBool(),
	}

	expenseShare, err := testQueries.CreateExpenseShare(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expenseShare)

	require.Equal(t, arg.ExpenseID, expenseShare.ExpenseID)
	require.Equal(t, arg.UserID, expenseShare.UserID)
	require.Equal(t, arg.Share, expenseShare.Share)
	require.Equal(t, arg.PaidStatus, expenseShare.PaidStatus)
	require.NotZero(t, expenseShare.ID)
	require.NotZero(t, expenseShare.CreatedAt)

	return expenseShare

}

func TestCreateExpenseShare(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	expense := createRandomExpense(t, group)
	createRandomExpenseShare(t, expense)
}

func TestGetExpenseShareByID(t *testing.T) {
	user := createRandomUser(t)

	group := createRandomGroup(t, user)
	expense := createRandomExpense(t, group)
	expenseShare1 := createRandomExpenseShare(t, expense)
	expenseShare2, err := testQueries.GetExpenseShareByID(context.Background(), expenseShare1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, expenseShare2)

	require.Equal(t, expenseShare1.ID, expenseShare2.ID)
	require.Equal(t, expenseShare1.ExpenseID, expenseShare2.ExpenseID)
	require.Equal(t, expenseShare1.UserID, expenseShare2.UserID)
	require.Equal(t, expenseShare1.Share, expenseShare2.Share)
	require.Equal(t, expenseShare1.PaidStatus, expenseShare2.PaidStatus)
	require.WithinDuration(t, expenseShare1.CreatedAt, expenseShare2.CreatedAt, time.Second)
}

func TestListExpenseShares(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	expense := createRandomExpense(t, group)
	for i := 0; i < 10; i++ {
		createRandomExpenseShare(t, expense)
	}

	arg := ListExpenseSharesParams{
		Limit:     5,
		Offset:    5,
		ExpenseID: expense.ID,
	}

	expenseShares, err := testQueries.ListExpenseShares(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, expenseShares, 5)

	for _, expenseShare := range expenseShares {
		require.NotEmpty(t, expenseShare)
	}
}

func TestDeleteExpenseShare(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	expense := createRandomExpense(t, group)
	expenseShare1 := createRandomExpenseShare(t, expense)
	err := testQueries.DeleteExpenseShare(context.Background(), expenseShare1.ID)
	require.NoError(t, err)

	expenseShare2, err := testQueries.GetExpenseShareByID(context.Background(), expenseShare1.ID)
	require.Error(t, err)
	require.Empty(t, expenseShare2)
}

func TestUpdateExpenseShare(t *testing.T) {
	user := createRandomUser(t)

	group := createRandomGroup(t, user)
	expense := createRandomExpense(t, group)
	expenseShare1 := createRandomExpenseShare(t, expense)

	arg := UpdateExpenseShareParams{
		ID:         expenseShare1.ID,
		Share:      util.RandomMoney(),
		PaidStatus: util.RandomBool(),
	}

	expenseShare2, err := testQueries.UpdateExpenseShare(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expenseShare2)

	require.Equal(t, expenseShare1.ID, expenseShare2.ID)
	require.Equal(t, expenseShare1.ExpenseID, expenseShare2.ExpenseID)
	require.Equal(t, expenseShare1.UserID, expenseShare2.UserID)
	require.Equal(t, arg.Share, expenseShare2.Share)
	require.Equal(t, arg.PaidStatus, expenseShare2.PaidStatus)
	require.WithinDuration(t, expenseShare1.CreatedAt, expenseShare2.CreatedAt, time.Second)
}
