package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type createTripRequest struct {
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

func (server *Server) createTrip(ctx *gin.Context) {
	var req createTripRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateTripParams{
		Title:     req.Title,
		StartDate: startDate,
		EndDate:   endDate,
		UserID:    authPayload.UserId,
	}
	trip, err := server.store.CreateTrip(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, trip)
}

type getTripRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getTrip(ctx *gin.Context) {
	var req getTripRequest
	if err := ctx.ShouldBindUri(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	trip, err := server.store.GetTrip(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, trip)
}
