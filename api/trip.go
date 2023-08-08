package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type createTripRequest struct {
	Title     string    `json:"title" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

func (server *Server) createTrip(ctx *gin.Context) {
	var req createTripRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateTripParams{
		Title:     req.Title,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		UserID:    authPayload.UserId,
	}
	trip, err := server.store.CreateTrip(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, trip)
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
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetTripParams{
		ID:     req.ID,
		UserID: authPayload.UserId,
	}

	trip, err := server.store.GetTrip(ctx, arg)
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

type listUserTripsRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=1,max=10"`
	Offset int32 `form:"offset" binding:"required,min=0"`
}

func (server *Server) listUserTrips(ctx *gin.Context) {
	var req listUserTripsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListTripParams{
		UserID: authPayload.UserId,
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}
	trips, err := server.store.ListTrip(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, trips)

}

type updateTripRequest struct {
	ID        int64     `json:"id" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

func (server *Server) updateTrip(ctx *gin.Context) {
	var req updateTripRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateTripParams{
		Title:     req.Title,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		ID:        req.ID,
		UserID:    authPayload.UserId,
	}

	log.Println(arg)
	trip, err := server.store.UpdateTrip(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, trip)
}

type deleteTripRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) deleteTrip(ctx *gin.Context) {
	var req deleteTripRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.DeleteTripParams{
		ID:     req.ID,
		UserID: authPayload.UserId,
	}
	err := server.store.DeleteTrip(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Trip deleted")
}
