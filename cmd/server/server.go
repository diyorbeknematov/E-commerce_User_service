package server

import (
	"auth-service/config"
	"auth-service/generated/user"
	"auth-service/service"
	"auth-service/storage/postgres"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Server interface {
	StartServer() error
}

type serverImpl struct {
	userRepo postgres.UserRepository
	logger   *slog.Logger
}

func NewServer(userRepo postgres.UserRepository, logger *slog.Logger) Server {
	return &serverImpl{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *serverImpl) StartServer() error {
	listen, err := net.Listen("tcp", config.Load().GRPC_USER_PORT)
	if err != nil {
		s.logger.Error("Error create to new listener", "error", err.Error())
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	service := service.NewUserService(s.userRepo, s.logger)

	user.RegisterUserServiceServer(server, service)

	s.logger.Info("server is running ", "PORT", config.Load().GRPC_USER_PORT)
	if err := server.Serve(listen); err != nil {
		s.logger.Error("Faild server is running", "error", err.Error())
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
