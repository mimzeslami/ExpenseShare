package db

import (
	"context"
	"testing"
	"time"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func createRandomGroup(t *testing.T, user Users) Groups {
	groupCategory := createRandomGroupCategory(t)

	arg := CreateGroupParams{
		Name:        util.RandomString(6),
		CategoryID:  groupCategory.ID,
		CreatedByID: user.ID,
	}

	group, err := testQueries.CreateGroup(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, group)
	require.Equal(t, arg.Name, group.Name)
	require.Equal(t, arg.CategoryID, group.CategoryID)
	require.Equal(t, arg.CreatedByID, group.CreatedByID)
	require.NotZero(t, group.ID)
	require.NotZero(t, group.CreatedAt)

	return group

}

func TestCreateGroup(t *testing.T) {
	user := createRandomUser(t)
	createRandomGroup(t, user)
}

func TestGetGroupById(t *testing.T) {
	user := createRandomUser(t)
	group1 := createRandomGroup(t, user)
	arg := GetGroupByIDParams{
		ID:          group1.ID,
		CreatedByID: group1.CreatedByID,
	}
	group2, err := testQueries.GetGroupByID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, group2)

	require.Equal(t, group1.ID, group2.ID)
	require.Equal(t, group1.Name, group2.Name)
	require.Equal(t, group1.CategoryID, group2.CategoryID)
	require.Equal(t, group1.CreatedByID, group2.CreatedByID)
	require.WithinDuration(t, group1.CreatedAt, group2.CreatedAt, time.Second)
}

func TestListGroup(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomGroup(t, user)
	}

	arg := ListGroupsParams{
		Limit:       5,
		Offset:      5,
		CreatedByID: user.ID,
	}

	groups, err := testQueries.ListGroups(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, groups, 5)

	for _, group := range groups {
		require.NotEmpty(t, group)
	}
}

func TestDeleteGroup(t *testing.T) {
	user := createRandomUser(t)
	group1 := createRandomGroup(t, user)
	arg := DeleteGroupParams{
		ID:          group1.ID,
		CreatedByID: group1.CreatedByID,
	}
	err := testQueries.DeleteGroup(context.Background(), arg)
	require.NoError(t, err)
	arg2 := GetGroupByIDParams{
		ID:          group1.ID,
		CreatedByID: group1.CreatedByID,
	}
	group2, err := testQueries.GetGroupByID(context.Background(), arg2)
	require.Error(t, err)
	require.Empty(t, group2)
}

func TestUpdateGroup(t *testing.T) {
	user := createRandomUser(t)
	group1 := createRandomGroup(t, user)

	arg := UpdateGroupParams{
		ID:          group1.ID,
		Name:        util.RandomString(6),
		CategoryID:  group1.CategoryID,
		CreatedByID: group1.CreatedByID,
	}

	group2, err := testQueries.UpdateGroup(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, group2)

	require.Equal(t, group1.ID, group2.ID)
	require.Equal(t, arg.Name, group2.Name)
	require.Equal(t, group1.CategoryID, group2.CategoryID)
	require.Equal(t, group1.CreatedByID, group2.CreatedByID)
	require.WithinDuration(t, group1.CreatedAt, group2.CreatedAt, time.Second)
}
