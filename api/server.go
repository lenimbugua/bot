package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/lenimbugua/artworld/db/sqlc"
	"github.com/lenimbugua/artworld/token"
	"github.com/lenimbugua/artworld/util"
)

// server serves HTTP  requests for artworld
type Server struct {
	config     util.Config
	dbStore    db.QueryStore
	tokenMaker token.Maker
	router     *gin.Engine
}

// Newserver creates a new HTTP server and sets up routing
func NewServer(config util.Config, dbStore db.QueryStore) (*Server, error) {
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
	server.router = router

}

// start runs the server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
