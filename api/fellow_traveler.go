package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
)

type createFellowTravelerRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	TripID    int64  `json:"trip_id" binding:"required"`
}

func (server *Server) createFellowTraveler(ctx *gin.Context) {
	var req createFellowTravelerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateFellowTravelersParams{
		FellowFirstName: req.FirstName,
		FellowLastName:  req.LastName,
		TripID:          req.TripID,
	}

	fellowTraveler, err := server.store.CreateFellowTravelers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, fellowTraveler)

}

type deleteFellowTravelerRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) deleteFellowTraveler(ctx *gin.Context) {
	var req deleteFellowTravelerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteFellowTraveler(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}

type updateFellowTravelerRequest struct {
	ID        int64  `json:"id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (server *Server) updateFellowTraveler(ctx *gin.Context) {
	var req updateFellowTravelerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateFellowTravelerParams{
		ID:              req.ID,
		FellowFirstName: req.FirstName,
		FellowLastName:  req.LastName,
	}
	fellowTraveler, err := server.store.UpdateFellowTraveler(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, fellowTraveler)
}

type listFellowTravelerRequest struct {
	TripID int64 `uri:"trip_id" binding:"required"`
}

func (server *Server) listFellowTraveler(ctx *gin.Context) {
	var req listFellowTravelerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetTripFellowTravelersParams{
		TripID: req.TripID,
		UserID: authPayload.UserId,
	}
	fellowTravelers, err := server.store.GetTripFellowTravelers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, fellowTravelers)
}

type getFellowTraveler struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getFellowTraveler(ctx *gin.Context) {
	var req getFellowTraveler
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetFellowTravelerParams{
		ID:     req.ID,
		UserID: authPayload.UserId,
	}
	fellowTraveler, err := server.store.GetFellowTraveler(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, fellowTraveler)
}
