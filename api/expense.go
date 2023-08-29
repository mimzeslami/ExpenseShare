package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type createExpenseRequest struct {
	GroupID     int64  `json:"group_id" binding:"required"`
	PaidByID    int64  `json:"paid_by_id" binding:"required"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

func (server *Server) createExpense(ctx *gin.Context) {
	var req createExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateExpenseParams{
		GroupID:     req.GroupID,
		PaidByID:    req.PaidByID,
		Amount:      req.Amount,
		Description: req.Description,
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if status, err := server.isUserGroupMember(ctx, authPayload.UserId, req.GroupID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else if status == false {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	expense, err := server.store.CreateExpense(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, expense)
}

type updateExpenseRequest struct {
	ID          int64  `json:"id" binding:"required"`
	GroupID     int64  `json:"group_id" binding:"required"`
	PaidByID    int64  `json:"paid_by_id" binding:"required"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

func (server *Server) updateExpense(ctx *gin.Context) {
	var req updateExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateExpenseParams{
		ID:          req.ID,
		Amount:      req.Amount,
		Description: req.Description,
	}

	expense, err := server.store.UpdateExpense(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, expense)
}

type listExpenseRequest struct {
	Limit   int32 `form:"limit" binding:"required,min=1,max=10"`
	Offset  int32 `form:"offset" binding:"required,min=0"`
	GroupID int64 `form:"group_id" binding:"required,min=1"`
}

func (server *Server) listExpenses(ctx *gin.Context) {
	var req listExpenseRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListExpensesParams{
		Limit:   req.Limit,
		Offset:  req.Offset,
		GroupID: req.GroupID,
	}

	expenses, err := server.store.ListExpenses(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, expenses)
}

type getExpenseByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getExpenseDetail(ctx *gin.Context) {
	var req getExpenseByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	expense, err := server.store.GetExpenseByID(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, expense)

}

type deleteExpenseRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteExpense(ctx *gin.Context) {
	var req deleteExpenseRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteExpense(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
