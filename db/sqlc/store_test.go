package db

import (
	"context"
	"testing"

	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func TestCreateGroupTx(t *testing.T) {

	store := NewStore(testDB)

	n := 10

	groupCategory := createRandomGroupCategory(t)
	user := createRandomUser(t)

	arg := CreateGroupParams{
		Name:        util.RandomString(6),
		CategoryID:  groupCategory.ID,
		CreatedByID: user.ID,
		ImagePath:   util.RandomString(6),
	}

	errs := make(chan error)
	results := make(chan CreateGroupTxResults)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.CreateGroupTx(context.Background(), arg)
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result.Group)
		require.NotZero(t, result.Group.ID)
		require.NotEmpty(t, result.GroupMembers)
		require.NotZero(t, result.GroupMembers.ID)
		require.Equal(t, arg.Name, result.Group.Name)
		require.Equal(t, arg.CategoryID, result.Group.CategoryID)
		require.Equal(t, result.Group.ID, result.GroupMembers.GroupID)
		require.Equal(t, arg.CreatedByID, result.Group.CreatedByID)

	}

}
func TestDeleteGroupTx(t *testing.T) {

	store := NewStore(testDB)

	n := 10

	groupCategory := createRandomGroupCategory(t)
	user := createRandomUser(t)

	arg := CreateGroupParams{

		Name:        util.RandomString(6),
		CategoryID:  groupCategory.ID,
		CreatedByID: user.ID,
		ImagePath:   util.RandomString(6),
	}
	group, err := store.CreateGroupTx(context.Background(), arg)

	require.NoError(t, err)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			err := store.DeleteGroupTx(context.Background(), DeleteGroupParams{
				ID:          group.Group.ID,
				CreatedByID: group.Group.CreatedByID,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

	}

}

func TestAddUserToGroupTx(t *testing.T) {
	store := NewStore(testDB)
	n := 10

	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	arg := AddUserToGroupParams{
		GroupID:      group.ID,
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Phone:        util.RandomString(6),
		Email:        util.RandomEmail(),
		TimeZone:     util.RandomString(6),
		GroupOwnerID: user.ID,
	}

	errs := make(chan error)
	results := make(chan AddUserToGroupResults)

	for i := 0; i < n; i++ {
		go func(i int) {
			arg.Phone = util.RandomString(9)
			result, err := store.AddUserToGroupTx(context.Background(), arg)
			errs <- err
			results <- result
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result.GroupMembers)
		require.NotZero(t, result.GroupMembers.ID)
		require.NotEmpty(t, result.User)
		require.NotZero(t, result.User.ID)
		require.NotEmpty(t, result.Invitations)
		require.NotZero(t, result.Invitations.ID)
		require.Equal(t, arg.GroupID, result.GroupMembers.GroupID)
		require.Equal(t, arg.GroupID, result.Invitations.GroupID)
		require.Equal(t, arg.GroupOwnerID, result.Invitations.InviterID)
	}
}
