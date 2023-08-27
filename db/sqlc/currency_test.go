package db

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mimzeslami/expense_share/util"

	"github.com/stretchr/testify/require"
)

func createRandomCurrency(t *testing.T) Currencies {

	arg := CreateCurrencyParams{
		Name:      util.RandomString(6),
		Code:      util.RandomString(3),
		Symbol:    util.RandomString(1),
		UpdatedAt: util.RandomDatetime(),
	}

	currency, err := testQueries.CreateCurrency(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, currency)

	require.Equal(t, arg.Name, currency.Name)
	require.NotZero(t, currency.ID)

	return currency
}

func TestCreateCurrency(t *testing.T) {
	createRandomCurrency(t)
}

func TestGetCurrency(t *testing.T) {
	currency1 := createRandomCurrency(t)
	currency2, err := testQueries.GetCurrencyByID(context.Background(), currency1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, currency2)

	require.Equal(t, currency1.ID, currency2.ID)
	require.Equal(t, currency1.Name, currency2.Name)
	require.Equal(t, currency1.Code, currency2.Code)
	require.Equal(t, currency1.Symbol, currency2.Symbol)
}

func TestListCurrencies(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCurrency(t)
	}

	arg := ListCurrenciesParams{
		Limit:  5,
		Offset: 5,
	}

	currencies, err := testQueries.ListCurrencies(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, currencies, 5)

	for _, currency := range currencies {
		require.NotEmpty(t, currency)
	}
}

func TestDeleteCurrency(t *testing.T) {
	currency1 := createRandomCurrency(t)
	err := testQueries.DeleteCurrency(context.Background(), currency1.ID)
	require.NoError(t, err)

	currency2, err := testQueries.GetCurrencyByID(context.Background(), currency1.ID)
	require.Error(t, err)
	require.Empty(t, currency2)
}
