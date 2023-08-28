package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomGroupMember(t *testing.T, group Groups) GroupMembers {
	user := createRandomUser(t)

	arg := CreateGroupMemberParams{
		GroupID: group.ID,
		UserID:  user.ID,
	}

	groupMember, err := testQueries.CreateGroupMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, groupMember)
	require.Equal(t, arg.GroupID, groupMember.GroupID)
	require.Equal(t, arg.UserID, groupMember.UserID)
	require.NotZero(t, groupMember.ID)
	require.NotZero(t, groupMember.CreatedAt)

	return groupMember

}

func TestCreateGroupMember(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	createRandomGroupMember(t, group)
}

func TestGetGroupMemberById(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	groupMember1 := createRandomGroupMember(t, group)
	groupMember2, err := testQueries.GetGroupMemberByID(context.Background(), groupMember1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, groupMember2)

	require.Equal(t, groupMember1.ID, groupMember2.ID)
	require.Equal(t, groupMember1.GroupID, groupMember2.GroupID)
	require.Equal(t, groupMember1.UserID, groupMember2.UserID)
	require.WithinDuration(t, groupMember1.CreatedAt, groupMember2.CreatedAt, time.Second)
}

func TestListGroupMembers(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	for i := 0; i < 10; i++ {
		createRandomGroupMember(t, group)
	}

	arg := ListGroupMembersParams{
		Limit:   5,
		Offset:  5,
		GroupID: group.ID,
	}

	groupMembers, err := testQueries.ListGroupMembers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, groupMembers, 5)

	for _, groupMember := range groupMembers {
		require.NotEmpty(t, groupMember)
	}
}

func TestDeleteGroupMember(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	groupMember1 := createRandomGroupMember(t, group)
	err := testQueries.DeleteGroupMember(context.Background(), groupMember1.ID)
	require.NoError(t, err)

	groupMember2, err := testQueries.GetGroupMemberByID(context.Background(), groupMember1.ID)
	require.Error(t, err)
	require.Empty(t, groupMember2)
}

func TestUpdateGroupMember(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	groupMember1 := createRandomGroupMember(t, group)

	arg := UpdateGroupMemberParams{
		ID:     groupMember1.ID,
		UserID: groupMember1.UserID,
	}

	groupMember2, err := testQueries.UpdateGroupMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, groupMember2)

	require.Equal(t, groupMember1.ID, groupMember2.ID)
	require.Equal(t, groupMember1.GroupID, groupMember2.GroupID)
	require.Equal(t, groupMember1.UserID, groupMember2.UserID)
	require.WithinDuration(t, groupMember1.CreatedAt, groupMember2.CreatedAt, time.Second)
}

func TestDeleteGroupMembersByGroupID(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	for i := 0; i < 10; i++ {
		createRandomGroupMember(t, group)
	}

	err := testQueries.DeleteGroupMembers(context.Background(), group.ID)
	require.NoError(t, err)

	groupMembers, err := testQueries.ListGroupMembers(context.Background(), ListGroupMembersParams{
		GroupID: group.ID,
	})
	require.NoError(t, err)
	require.Len(t, groupMembers, 0)
}

func TestGetGroupMembersWithDetail(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t, user)
	for i := 0; i < 10; i++ {
		createRandomGroupMember(t, group)
	}
	arg := ListGroupMembersWithDetailsParams{
		Limit:   5,
		Offset:  5,
		GroupID: group.ID,
	}

	groupMembers, err := testQueries.ListGroupMembersWithDetails(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, groupMembers, 5)

	for _, groupMember := range groupMembers {
		require.NotEmpty(t, groupMember)
	}
}
