package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lenimbugua/bot/token"
	"github.com/lenimbugua/bot/util"
)

// server serves HTTP  requests for bot
type Server struct {
	config     util.Config
	dbStore    db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// Newserver creates a new HTTP server and sets up routing
func NewServer(config util.Config, dbStore db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot Create token %w", err)
	}
	server := &Server{
		dbStore:    dbStore,
		tokenMaker: tokenMaker,
		config:     config,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/channels", server.createChannel)

	authRoutes.POST("/companies", server.createCompany)

	//getCompanyByEmail uses query string
	authRoutes.GET("/companies", server.getCompanyByEmail)
	//getCompanyByID uses query uri
	authRoutes.GET("/companies/:id", server.getCompanyByID)

	authRoutes.GET("/list/companies", server.listCompanies)

	authRoutes.PUT("/companies/:id", server.updateCompany)

	authRoutes.DELETE("/companies/:id", server.deleteCompany)

	authRoutes.POST("/bots", server.createBot)
	authRoutes.PUT("/bots/:id", server.updateBot)
	authRoutes.DELETE("/bots/:id", server.deleteBot)
	authRoutes.GET("/list/bots", server.listBots)
	authRoutes.GET("/list/companybots", server.listCompanyBots)

	server.router = router

}

// start runs the server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// create a basic gin error
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func descriptiveError(err string) gin.H {
	return gin.H{"error": err}
}
