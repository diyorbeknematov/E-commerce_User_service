package service

import (
	"auth-service/storage/postgres"
	"log/slog"
)

type ServiceManager interface {
	User() UserService
	Auth() AuthService
}

type serviceManager struct {
	user        UserService
	authService AuthService
}

func NewServiceManager(storage postgres.MainRepository, logger *slog.Logger) ServiceManager {
	return &serviceManager{
		user:        NewUserService(storage.Users(), logger),
		authService: NewAuthenticatonService(storage.Authentications(), logger),
	}
}

func (s *serviceManager) User() UserService {
	return s.user
}

func (s *serviceManager) Auth() AuthService {
	return s.authService
}
