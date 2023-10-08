package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/util"
)

type createUserRequest struct {
	FirstName string `json:"first_name" binding:"required,alphanum"`
	LastName  string `json:"last_name" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	ImagePath string `json:"image_path"`
	TimeZone  string `json:"time_zone"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ImagePath string    `json:"image_path"`
	TimeZone  string    `json:"time_zone"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.Users) userResponse {
	return userResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		TimeZone:  user.TimeZone,
		ImagePath: user.ImagePath,
		CreatedAt: user.CreatedAt,
	}
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	User                 userResponse `json:"user"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		TimeZone:     req.TimeZone,
		ImagePath:    req.ImagePath,
		PasswordHash: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == db.UniqueViolation {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("user already exists with this email")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) login(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid email or password")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.CheckPassword(req.Password, user.PasswordHash); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid email or password")))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{

		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		User:                 newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)

}

type updateUserRequest struct {
	ID             int64  `json:"id" binding:"required,min=1"`
	FirstName      string `json:"first_name" binding:"required,alphanum"`
	LastName       string `json:"last_name" binding:"required,alphanum"`
	Password       string `json:"password" binding:"required,min=6"`
	Email          string `json:"email" binding:"email"`
	Phone          string `json:"phone" binding:"required"`
	ImagePath      string `json:"image_path"`
	TimeZone       string `json:"time_zone"`
	InvitationCode string `json:"code"`
}

func (server *Server) completeProfile(ctx *gin.Context) {
	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	invitation, err := server.store.GetInvitationByCode(ctx, req.InvitationCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if invitation.Status != db.InvitationStatusPending {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invitation is not pending")))
		return
	}
	arg := db.UpdateInvitationParams{
		ID:     invitation.ID,
		Status: db.InvitationStatusAccepted,
		AcceptedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	invitation, err = server.store.UpdateInvitation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userInfo := db.UpdateUserParams{
		ID:           req.ID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		TimeZone:     req.TimeZone,
		ImagePath:    req.ImagePath,
		PasswordHash: hashedPassword,
	}

	user, err := server.store.UpdateUser(ctx, userInfo)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == db.UniqueViolation {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("user already exists with this email")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) googleLoginRequest(ctx *gin.Context) {
	state := util.RandomString(12)
	config, err := util.NewOAuthConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, config.AuthCodeURL(state))
}

func (server *Server) googleLoginCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	config, err := util.NewOAuthConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Println("code: ", code)
	log.Println("config: ", config.ClientSecret)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "could not exchange oauth code")
		return
	}

	if !token.Valid() {
		ctx.JSON(http.StatusInternalServerError, "invalid access token")
		return
	}

	userInfo, err := getUserInfoFromGoogle(token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userInfo)

}

type GoogleUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	ImageURL  string `json:"picture"`
}

func getUserInfoFromGoogle(token string) (GoogleUserInfo, error) {
	var userInfo GoogleUserInfo
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return userInfo, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	// Read the response body into a []byte variable
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return userInfo, err
	}

	// Unmarshal the JSON from the []byte
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return userInfo, err
	}
	return userInfo, nil
}
