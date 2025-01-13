package api

import (
	"fmt"

	db "github.com/emzhofb/todo-crud/db/sqlc"
	"github.com/emzhofb/todo-crud/token"
	"github.com/emzhofb/todo-crud/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	connPool   *pgxpool.Pool
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store, connPool *pgxpool.Pool) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		connPool:   connPool,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("status", validStatus)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	router.POST("/tasks", server.createTask)
	router.GET("/tasks", server.listTask)
	router.GET("/tasks/:id", server.getTask)
	router.PUT("/tasks/:id", server.updateTask)
	router.DELETE("/tasks/:id", server.deleteTask)

	// use auth middleware
	router.GET("/token", server.createToken)
	authRoutes := router.Group("/api/v1").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/tasks", server.createTask)
	authRoutes.GET("/tasks", server.listTask)
	authRoutes.GET("/tasks/:id", server.getTask)
	authRoutes.PUT("/tasks/:id", server.updateTask)
	authRoutes.DELETE("/tasks/:id", server.deleteTask)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
