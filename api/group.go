package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type createGroupRequest struct {
	Name       string `json:"name" binding:"required"`
	CategoryID int64  `json:"category_id" binding:"required"`
	ImagePath  string `json:"image_path"`
}

type groupResponse struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	CategoryID  int64           `json:"category_id"`
	ImagePath   string          `json:"image_path"`
	CreatedByID int64           `json:"created_by_id"`
	GroupMember db.GroupMembers `json:"group_members"`
}

func newGroupResponse(group db.CreateGroupTxResults) groupResponse {
	return groupResponse{
		ID:          group.Group.ID,
		Name:        group.Group.Name,
		CategoryID:  group.Group.CategoryID,
		ImagePath:   group.Group.ImagePath,
		CreatedByID: group.Group.CreatedByID,
		GroupMember: group.GroupMembers,
	}
}

func (server *Server) createGroup(ctx *gin.Context) {
	var req createGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateGroupParams{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		CreatedByID: authPayload.UserId,
		ImagePath:   req.ImagePath,
	}

	group, err := server.store.CreateGroupTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, newGroupResponse(group))
}

type listGroupRequest struct {
	Offset int32 `form:"offset" binding:"required,min=0"`
	Limit  int32 `form:"limit" binding:"required,min=1,max=10"`
}

func (server *Server) listGroups(ctx *gin.Context) {
	var req listGroupRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListGroupsParams{
		Limit:       req.Limit,
		Offset:      (req.Offset - 1) * req.Limit,
		CreatedByID: authPayload.UserId,
	}

	groups, err := server.store.ListGroups(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, groups)
}

type getGroupRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getGroup(ctx *gin.Context) {
	var req getGroupRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetGroupByIDParams{
		ID:          req.ID,
		CreatedByID: authPayload.UserId,
	}

	group, err := server.store.GetGroupByID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, group)
}

type updateGroupRequest struct {
	ID         int64  `json:"id" binding:"required,min=1"`
	Name       string `json:"name" binding:"required"`
	CategoryID int64  `json:"category_id" binding:"required"`
	ImagePath  string `json:"image_path"`
}

func (server *Server) updateGroup(ctx *gin.Context) {
	var req updateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if ok, err := server.isUserGroupOwner(ctx, authPayload.UserId, req.ID); !ok {
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.UpdateGroupParams{
		ID:          req.ID,
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		CreatedByID: authPayload.UserId,
		ImagePath:   req.ImagePath,
	}

	group, err := server.store.UpdateGroup(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, group)
}

type deleteGroupRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteGroup(ctx *gin.Context) {
	var req deleteGroupRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if ok, err := server.isUserGroupOwner(ctx, authPayload.UserId, req.ID); !ok {
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.DeleteGroupParams{
		ID:          req.ID,
		CreatedByID: authPayload.UserId,
	}
	err := server.store.DeleteGroupTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
