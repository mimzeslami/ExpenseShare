package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomExpense(t *testing.T, user Users) Expenses {
	traveler := createRandomFellowTraveller(t, user)
	arg := CreateExpenseParams{
		TripID:          traveler.TripID,
		PayerTravelerID: traveler.ID,
		Amount:          util.RandomMoney(),
		Description:     util.RandomString(6),
	}
	expense, err := testQueries.CreateExpense(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expense)
	require.Equal(t, arg.TripID, expense.TripID)
	require.Equal(t, arg.PayerTravelerID, expense.PayerTravelerID)
	require.Equal(t, arg.Amount, expense.Amount)
	require.Equal(t, arg.Description, expense.Description)
	require.NotZero(t, expense.ID)

	return expense

}
func TestCreateExpense(t *testing.T) {
	user := createRandomUser(t)
	createRandomExpense(t, user)
}

func TestGetExpense(t *testing.T) {
	user := createRandomUser(t)

	expense1 := createRandomExpense(t, user)
	arg := GetExpenseParams{
		ID:     expense1.ID,
		UserID: user.ID,
	}
	expense2, err := testQueries.GetExpense(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expense2)
	require.Equal(t, expense1.ID, expense2.ID)
	require.Equal(t, expense1.TripID, expense2.TripID)
	require.Equal(t, expense1.PayerTravelerID, expense2.PayerTravelerID)
	require.Equal(t, expense1.Amount, expense2.Amount)
	require.Equal(t, expense1.Description, expense2.Description)
}

func TestGetTripExpenses(t *testing.T) {
	user := createRandomUser(t)

	expense1 := createRandomExpense(t, user)
	arg := GetTripExpensesParams{
		TripID: expense1.TripID,
		UserID: user.ID,
	}
	expense2, err := testQueries.GetTripExpenses(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expense2)
	require.Equal(t, expense1.ID, expense2[0].ID)
	require.Equal(t, expense1.TripID, expense2[0].TripID)
	require.Equal(t, expense1.PayerTravelerID, expense2[0].PayerTravelerID)
	require.Equal(t, expense1.Amount, expense2[0].Amount)
	require.Equal(t, expense1.Description, expense2[0].Description)
}

func TestDeleteExpense(t *testing.T) {
	user := createRandomUser(t)
	expense1 := createRandomExpense(t, user)
	err := testQueries.DeleteExpense(context.Background(), expense1.ID)
	require.NoError(t, err)
	expense2, err := testQueries.GetExpense(context.Background(), GetExpenseParams{
		ID:     expense1.ID,
		UserID: user.ID,
	})
	require.Error(t, err)
	require.Empty(t, expense2)
	require.ErrorIs(t, err, sql.ErrNoRows)

}

func TestDeleteTripExpensees(t *testing.T) {
	user := createRandomUser(t)
	expense1 := createRandomExpense(t, user)
	err := testQueries.DeleteTripExpenses(context.Background(), expense1.TripID)
	require.NoError(t, err)

	expenses, err := testQueries.GetTripExpenses(context.Background(), GetTripExpensesParams{
		TripID: expense1.TripID,
		UserID: user.ID,
	})
	require.Empty(t, expenses)

}

func TestUpdateExpense(t *testing.T) {
	user := createRandomUser(t)
	expense1 := createRandomExpense(t, user)
	arg := UpdateExpenseParams{
		ID:              expense1.ID,
		TripID:          expense1.TripID,
		PayerTravelerID: expense1.PayerTravelerID,
		Amount:          util.RandomMoney(),
		Description:     util.RandomString(6),
	}
	expense2, err := testQueries.UpdateExpense(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expense2)
	require.Equal(t, expense1.ID, expense2.ID)
	require.Equal(t, arg.Amount, expense2.Amount)
	require.Equal(t, arg.Description, expense2.Description)
}
