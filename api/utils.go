package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
)

func (server *Server) isUserGroupMember(ctx *gin.Context, userID int64, groupID int64) (status bool, err error) {
	arg := db.GetGroupMemberByGroupIDAndUserIDParams{
		GroupID: groupID,
		UserID:  userID,
	}

	group_member, err := server.store.GetGroupMemberByGroupIDAndUserID(ctx, arg)
	if err != nil {
		return false, err
	}

	if group_member.GroupID == groupID && group_member.UserID == userID {
		return true, nil
	}
	return false, nil

}

func (server *Server) isUserGroupOwner(ctx *gin.Context, userID int64, groupID int64) (status bool, err error) {
	arg := db.GetGroupByGroupIDAndUserIDParams{
		ID:          groupID,
		CreatedByID: userID,
	}

	group, err := server.store.GetGroupByGroupIDAndUserID(ctx, arg)
	if err != nil {
		return false, err
	}

	if group.ID == groupID && group.CreatedByID == userID {
		return true, nil
	}
	return false, nil

}
