package db

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/mimzeslami/expense_share/util"

	"github.com/stretchr/testify/require"
)

func createRandomInvitation(t *testing.T, group Groups, user Users) Invitations {
	user2 := createRandomUser(t)

	arg := CreateInvitationParams{
		GroupID:   group.ID,
		InviterID: user2.ID,
		InviteeID: user.ID,
		Code:      util.RandomString(6),
		Status:    "pending",
	}

	invitation, err := testQueries.CreateInvitation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invitation)
	require.Equal(t, arg.GroupID, invitation.GroupID)
	require.Equal(t, arg.InviterID, invitation.InviterID)
	require.Equal(t, arg.InviteeID, invitation.InviteeID)
	require.Equal(t, arg.Code, invitation.Code)
	require.NotZero(t, invitation.ID)
	require.NotZero(t, invitation.CreatedAt)

	return invitation

}

func TestCreateInvitation(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t)
	createRandomInvitation(t, group, user)
}

func TestGetInvitationByID(t *testing.T) {
	group := createRandomGroup(t)
	user := createRandomUser(t)
	invitation1 := createRandomInvitation(t, group, user)
	invitation2, err := testQueries.GetInvitationByID(context.Background(), invitation1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, invitation2)

	require.Equal(t, invitation1.ID, invitation2.ID)
	require.Equal(t, invitation1.GroupID, invitation2.GroupID)
	require.Equal(t, invitation1.InviterID, invitation2.InviterID)
	require.Equal(t, invitation1.InviteeID, invitation2.InviteeID)
	require.Equal(t, invitation1.Code, invitation2.Code)
	require.Equal(t, invitation1.Status, invitation2.Status)
	require.WithinDuration(t, invitation1.CreatedAt, invitation2.CreatedAt, time.Second)
}

func TestListInvitations(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t)
	for i := 0; i < 10; i++ {
		createRandomInvitation(t, group, user)
	}

	arg := ListInvitationsForInviteeParams{
		Limit:     5,
		Offset:    5,
		InviteeID: user.ID,
	}

	invitations, err := testQueries.ListInvitationsForInvitee(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, invitations, 5)

	for _, invitation := range invitations {
		require.NotEmpty(t, invitation)
	}
}

func TestDeleteInvitation(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t)
	invitation := createRandomInvitation(t, group, user)
	err := testQueries.DeleteInvitation(context.Background(), invitation.ID)
	require.NoError(t, err)
}

func TestUpdateInvitation(t *testing.T) {
	user := createRandomUser(t)
	group := createRandomGroup(t)
	invitation1 := createRandomInvitation(t, group, user)

	arg := UpdateInvitationParams{
		ID:     invitation1.ID,
		Status: "accepted",
	}

	invitation2, err := testQueries.UpdateInvitation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invitation2)

	require.Equal(t, invitation1.ID, invitation2.ID)
	require.Equal(t, invitation1.GroupID, invitation2.GroupID)
	require.Equal(t, invitation1.InviterID, invitation2.InviterID)
	require.Equal(t, invitation1.InviteeID, invitation2.InviteeID)
	require.Equal(t, invitation1.Code, invitation2.Code)
	require.Equal(t, arg.Status, invitation2.Status)
	require.WithinDuration(t, invitation1.CreatedAt, invitation2.CreatedAt, time.Second)
}
