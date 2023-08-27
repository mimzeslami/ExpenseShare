package db

import (
	"context"
	"testing"
	"time"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomExpense(t *testing.T, group Groups) Expenses {
	user := createRandomUser(t)

	arg := CreateExpenseParams{
		GroupID:     group.ID,
		PaidByID:    user.ID,
		Description: util.RandomString(6),
		Date:        util.RandomDatetime(),
		Amount:      util.RandomMoney(),
	}

	expense, err := testQueries.CreateExpense(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expense)
	require.Equal(t, arg.GroupID, expense.GroupID)
	require.Equal(t, arg.PaidByID, expense.PaidByID)
	require.Equal(t, arg.Description, expense.Description)
	require.Equal(t, arg.Date, expense.Date)
	require.Equal(t, arg.Amount, expense.Amount)
	require.NotZero(t, expense.ID)
	require.NotZero(t, expense.CreatedAt)

	return expense

}

func TestCreateExpense(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	createRandomExpense(t, group)
}

func TestGetExpenseByID(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	expense1 := createRandomExpense(t, group)
	expense2, err := testQueries.GetExpenseByID(context.Background(), expense1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, expense2)

	require.Equal(t, expense1.ID, expense2.ID)
	require.Equal(t, expense1.GroupID, expense2.GroupID)
	require.Equal(t, expense1.PaidByID, expense2.PaidByID)
	require.Equal(t, expense1.Description, expense2.Description)
	require.Equal(t, expense1.Date, expense2.Date)
	require.Equal(t, expense1.Amount, expense2.Amount)
	require.WithinDuration(t, expense1.CreatedAt, expense2.CreatedAt, time.Second)
}

func TestListExpense(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	for i := 0; i < 10; i++ {
		createRandomExpense(t, group)
	}

	arg := ListExpensesParams{
		Limit:   5,
		Offset:  5,
		GroupID: group.ID,
	}

	expenses, err := testQueries.ListExpenses(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, expenses, 5)

	for _, expense := range expenses {
		require.NotEmpty(t, expense)
	}
}

func TestDeleteExpense(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	expense1 := createRandomExpense(t, group)
	err := testQueries.DeleteExpense(context.Background(), expense1.ID)
	require.NoError(t, err)

	expense2, err := testQueries.GetExpenseByID(context.Background(), expense1.ID)
	require.Error(t, err)
	require.Empty(t, expense2)
}

func TestUpdateExpense(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	expense1 := createRandomExpense(t, group)

	arg := UpdateExpenseParams{
		ID:          expense1.ID,
		Description: util.RandomString(6),
		Date:        util.RandomDatetime(),
		Amount:      util.RandomMoney(),
	}

	expense2, err := testQueries.UpdateExpense(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expense2)
	require.Equal(t, expense1.ID, expense2.ID)
	require.Equal(t, expense1.GroupID, expense2.GroupID)
	require.Equal(t, expense1.PaidByID, expense2.PaidByID)
	require.Equal(t, arg.Description, expense2.Description)
	require.Equal(t, arg.Date, expense2.Date)
	require.Equal(t, arg.Amount, expense2.Amount)
	require.WithinDuration(t, expense1.CreatedAt, expense2.CreatedAt, time.Second)
}
