package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
	"github.com/mimzeslami/expense_share/util"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}
func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.login)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/group_categories", server.createGroupCategory)
	authRoutes.GET("/group_categories", server.listGroupCategories)
	authRoutes.GET("/group_categories/:id", server.getGroupCategory)
	authRoutes.PUT("/group_categories", server.updateGroupCategory)
	authRoutes.DELETE("/group_categories/:id", server.deleteGroupCategory)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
