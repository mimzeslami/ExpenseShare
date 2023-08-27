package db

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/mimzeslami/expense_share/util"

	"github.com/stretchr/testify/require"
)

func createRandomGroupCategory(t *testing.T) GroupCategories {

	categoryName := util.RandomString(6)

	groupCategory, err := testQueries.CreateGroupCategory(context.Background(), categoryName)
	require.NoError(t, err)
	require.NotEmpty(t, groupCategory)

	require.Equal(t, categoryName, groupCategory.Name)
	require.NotZero(t, groupCategory.ID)
	require.NotZero(t, groupCategory.CreatedAt)

	return groupCategory
}

func TestCreateGroupCategory(t *testing.T) {
	createRandomGroupCategory(t)
}

func TestGetGroupCategory(t *testing.T) {
	groupCategory1 := createRandomGroupCategory(t)
	groupCategory2, err := testQueries.GetGroupCategory(context.Background(), groupCategory1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, groupCategory2)

	require.Equal(t, groupCategory1.ID, groupCategory2.ID)
	require.Equal(t, groupCategory1.Name, groupCategory2.Name)
	require.WithinDuration(t, groupCategory1.CreatedAt, groupCategory2.CreatedAt, time.Second)
}

func TestListGroupCategories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomGroupCategory(t)
	}

	arg := ListGroupCategoriesParams{
		Limit:  5,
		Offset: 5,
	}

	groupCategories, err := testQueries.ListGroupCategories(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, groupCategories, 5)

	for _, groupCategory := range groupCategories {
		require.NotEmpty(t, groupCategory)
	}
}

func TestDeleteGroupCategory(t *testing.T) {
	groupCategory1 := createRandomGroupCategory(t)
	err := testQueries.DeleteGroupCategory(context.Background(), groupCategory1.ID)
	require.NoError(t, err)

	groupCategory2, err := testQueries.GetGroupCategory(context.Background(), groupCategory1.ID)
	require.Error(t, err)
	require.Empty(t, groupCategory2)
}

func TestUpdateGroupCategory(t *testing.T) {
	groupCategory1 := createRandomGroupCategory(t)

	arg := UpdateGroupCategoryParams{
		ID:   groupCategory1.ID,
		Name: util.RandomString(6),
	}

	groupCategory2, err := testQueries.UpdateGroupCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, groupCategory2)

	require.Equal(t, groupCategory1.ID, groupCategory2.ID)
	require.NotEqual(t, groupCategory1.Name, groupCategory2.Name)
	require.WithinDuration(t, groupCategory1.CreatedAt, groupCategory2.CreatedAt, time.Second)
}
