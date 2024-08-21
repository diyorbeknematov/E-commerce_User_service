package api

import (
	"auth-service/api/handler"
	_ "auth-service/api/handler/docs"
	"auth-service/api/middleware"
	"auth-service/config"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller interface {
	SetupRoutes(handler.AutheticationHandler, *slog.Logger)
	StartServer(config.Config) error
}

type controllerImpl struct {
	Port   string
	router *gin.Engine
}

func NewController(router *gin.Engine) Controller {
	return &controllerImpl{router: router}
}

func (c *controllerImpl) StartServer(cfg config.Config) error {
	if c.Port == "" {
		log.Println("Server port is not set, using default port 8081")
		c.Port = cfg.HTTP_PORT
	}
	return c.router.Run(c.Port)
}

// @title Authentication service
// @version 1.0
// @description service for authentication user
// @host localhost:8081
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api
// @schemes http
func (c *controllerImpl) SetupRoutes(handler handler.AutheticationHandler, logger *slog.Logger) {
	// Implement routes setup here

	// Swagger endpointini sozlash
	c.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := c.router.Group("/api")
	router.Use(middleware.LogMiddleware(logger))
	{
		router.POST("/register", handler.Register)
		router.POST("/login", handler.Login)
		router.POST("/refresh", handler.Refresh)
		router.POST("/logout", middleware.IsAuthenticated(), handler.Logout)
	}
}
