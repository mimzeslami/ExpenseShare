package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
)

type createGroupCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type groupCategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func newGroupCategoryResponse(groupCategory db.GroupCategories) groupCategoryResponse {
	return groupCategoryResponse{
		Name: groupCategory.Name,
		ID:   groupCategory.ID,
	}
}

func (server *Server) createGroupCategory(ctx *gin.Context) {
	var req createGroupCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	groupCategory, err := server.store.CreateGroupCategory(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, newGroupCategoryResponse(groupCategory))
}

type listGroupCategoriesRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=1,max=10"`
	Offset int32 `form:"offset" binding:"required,min=0"`
}

func (server *Server) listGroupCategories(ctx *gin.Context) {
	var req listGroupCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListGroupCategoriesParams{
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}
	groupCategories, err := server.store.ListGroupCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, groupCategories)
}

type getGroupCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getGroupCategory(ctx *gin.Context) {
	var req getGroupCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	groupCategory, err := server.store.GetGroupCategory(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, groupCategory)
}

type updateGroupCategoryRequest struct {
	ID   int64  `json:"id" binding:"required,min=1"`
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateGroupCategory(ctx *gin.Context) {
	var req updateGroupCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	groupCategory, err := server.store.UpdateGroupCategory(ctx, db.UpdateGroupCategoryParams{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, groupCategory)
}

type deleteGroupCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteGroupCategory(ctx *gin.Context) {
	var req deleteGroupCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteGroupCategory(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
