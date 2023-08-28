package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type getCurrentInvitationStatusRequest struct {
	GroupID   int64 `uri:"group_id" binding:"required,min=1"`
	InviteeID int64 `uri:"invitee_id" binding:"required,min=1"`
}

func (server *Server) getCurrentInvitationStatus(ctx *gin.Context) {
	var req getCurrentInvitationStatusRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetCurrentInvitationByGroupIDAndInviteeIDParams{
		GroupID:   req.GroupID,
		InviteeID: req.InviteeID,
		InviterID: authPayload.UserId,
	}

	invitation, err := server.store.GetCurrentInvitationByGroupIDAndInviteeID(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, invitation)

}

type getUserInfoByInvitationCodeRequest struct {
	InvitationCode string `uri:"code" binding:"required,min=1"`
}

func (server *Server) getUserInfoByInvitationCode(ctx *gin.Context) {
	var req getUserInfoByInvitationCodeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	invitation, err := server.store.GetInvitationByCode(ctx, req.InvitationCode)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if invitation.Status != db.InvitationStatusPending {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("invitation is not pending")))
		return
	}

	user, err := server.store.GetUserByID(ctx, invitation.InviteeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)

}

type createInvitationRequest struct {
	GroupID   int64 `json:"group_id" binding:"required,min=1"`
	InviteeID int64 `json:"invitee_id" binding:"required,min=1"`
}

func (server *Server) createInvitation(ctx *gin.Context) {
	var req createInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateInvitationParams{
		GroupID:   req.GroupID,
		InviterID: authPayload.UserId,
		InviteeID: req.InviteeID,
		Status:    db.InvitationStatusPending,
	}

	invitation, err := server.store.CreateInvitation(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, invitation)

}
