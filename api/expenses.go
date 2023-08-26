package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type createExpenseRequest struct {
	TripID          int64  `json:"trip_id" binding:"required"`
	PayerTravelerID int64  `json:"payer_traveler_id" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	Description     string `json:"description" binding:"required"`
}

func (server *Server) createExpense(ctx *gin.Context) {
	var req createExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateExpenseParams{
		TripID:          req.TripID,
		PayerTravelerID: req.PayerTravelerID,
		Amount:          req.Amount,
		Description:     req.Description,
	}

	expense, err := server.store.CreateExpense(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, expense)

}

type getExpansesRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getExpanses(ctx *gin.Context) {
	var req getExpansesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetExpenseParams{
		ID:     req.ID,
		UserID: authPayload.UserId,
	}
	expanses, err := server.store.GetExpense(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, expanses)
}

type getTripExpensesRequest struct {
	TripID int64 `uri:"trip_id" binding:"required"`
}

func (server *Server) getTripExpenses(ctx *gin.Context) {
	var req getTripExpensesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetTripExpensesParams{
		TripID: req.TripID,
		UserID: authPayload.UserId,
	}
	expenses, err := server.store.GetTripExpenses(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, expenses)
}

type updateExpenseRequest struct {
	ID              int64  `json:"id" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	Description     string `json:"description" binding:"required"`
	TripID          int64  `json:"trip_id" binding:"required"`
	PayerTravelerID int64  `json:"payer_traveler_id" binding:"required"`
}

func (server *Server) updateExpense(ctx *gin.Context) {
	var req updateExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateExpenseParams{
		ID:              req.ID,
		Amount:          req.Amount,
		Description:     req.Description,
		TripID:          req.TripID,
		PayerTravelerID: req.PayerTravelerID,
	}
	expense, err := server.store.UpdateExpense(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, expense)
}

type deleteExpenseRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (server *Server) deleteExpense(ctx *gin.Context) {
	var req deleteExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteExpense(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Expense deleted successfully")
}

type deleteTripExpansesRequest struct {
	trip_id int64 `uri:"trip_id" binding:"required"`
}

func (server *Server) deleteTripExpanses(ctx *gin.Context) {
	var req deleteTripExpansesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteTripExpenses(ctx, req.trip_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Expenses deleted successfully")
}
