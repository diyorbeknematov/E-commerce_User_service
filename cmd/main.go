package main

import (
	"auth-service/api"
	"auth-service/api/handler"
	"auth-service/cmd/server"
	"auth-service/config"
	"auth-service/logs"
	"auth-service/service"
	"auth-service/storage/postgres"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logs.InitLogger()
	logger.Info("Starting auth-service...")

	db, err := postgres.Connection()
	if err != nil {
		logger.Error("Error connecting to database", "error", err.Error())
		log.Fatal(err)
	}
	defer db.Close()

	serviceManager := service.NewServiceManager(postgres.NewMainRepository(db), logger)

	go func() {
		server := server.NewServer(postgres.NewUserRepository(db), logger)

		log.Printf("Starting gRPC server...%s", config.Load().GRPC_USER_PORT)
		log.Fatal(server.StartServer())
	}()

	h := handler.NewAutheticationHandler(serviceManager, logger)

	controller := api.NewController(gin.Default())
	controller.SetupRoutes(h, logger)

	controller.StartServer(config.Load())
}
