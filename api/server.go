package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/katatrina/my-simple-bank/applayer"
	"github.com/katatrina/my-simple-bank/token"
	"github.com/katatrina/my-simple-bank/util"
)

// HTTPServer serves HTTP requests for our banking service.
type HTTPServer struct {
	config     util.Config
	app        applayer.App
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewHTTPServer creates a new HTTP server and sets up routing.
func NewHTTPServer(app applayer.App, config util.Config, tokenMaker token.Maker) (*HTTPServer, error) {
	server := &HTTPServer{
		config:     config,
		app:        app,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *HTTPServer) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.makeTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *HTTPServer) Start(address string) error {
	return server.router.Run(address)
}
