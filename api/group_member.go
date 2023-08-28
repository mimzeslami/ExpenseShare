package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type createGroupMemberRequest struct {
	GroupID   int64  `json:"group_id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	TimeZone  string `json:"time_zone"`
}
type groupMemberResponse struct {
	User        userResponse    `json:"user"`
	GroupMember db.GroupMembers `json:"group_member"`
	Invitation  db.Invitations  `json:"invitation"`
}

func createGroupMemberResponse(result db.AddUserToGroupResults) groupMemberResponse {
	return groupMemberResponse{
		User: userResponse{
			FirstName: result.User.FirstName,
			LastName:  result.User.LastName,
			Email:     result.User.Email,
			Phone:     result.User.Phone,
			ImagePath: result.User.ImagePath,
			TimeZone:  result.User.TimeZone,
			CreatedAt: result.User.CreatedAt,
		},
		GroupMember: result.GroupMembers,
		Invitation:  result.Invitations,
	}
}

func (server *Server) createGroupMember(ctx *gin.Context) {
	var req createGroupMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.AddUserToGroupParams{
		GroupID:      req.GroupID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Email:        req.Email,
		TimeZone:     req.TimeZone,
		GroupOwnerID: authPayload.UserId,
	}

	groupMember, err := server.store.AddUserToGroupTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// send email or sms to user

	ctx.JSON(http.StatusOK, createGroupMemberResponse(groupMember))
}

type listGroupMembersRequest struct {
	Offset  int32 `form:"offset" binding:"required,min=0"`
	Limit   int32 `form:"limit" binding:"required,min=1,max=10"`
	GroupID int64 `form:"group_id" binding:"required"`
}

func (server *Server) listGroupMembers(ctx *gin.Context) {
	var req listGroupMembersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListGroupMembersWithDetailsParams{
		GroupID: req.GroupID,
		Limit:   req.Limit,
		Offset:  (req.Offset - 1) * req.Limit,
	}

	groupMembers, err := server.store.ListGroupMembersWithDetails(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, groupMembers)

}

type deleteGroupMemberRequest struct {
	ID      int64 `uri:"id" binding:"required,min=1"`
	GroupID int64 `uri:"group_id" binding:"required,min=1"`
}

func (server *Server) deleteGroupMember(ctx *gin.Context) {
	var req deleteGroupMemberRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	groupMember, err := server.store.GetGroupMemberByID(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if groupMember.GroupID != req.GroupID {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	err = server.store.DeleteGroupMember(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, req.ID)

}
