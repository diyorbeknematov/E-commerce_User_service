package service

import (
	"auth-service/api/token"
	"auth-service/models"
	"auth-service/storage/postgres"
	"context"
	"log/slog"
)

type AuthService interface {
	Register(ctx context.Context, user models.Register) (*models.Response, error)
	Login(ctx context.Context, login models.LoginRequest) (*models.User, error)
	SaveRefreshToken(ctx context.Context, req models.TokenRequest) error
	InvalidateRefreshToken(ctx context.Context, username string) error
	IsRefreshTokenValid(ctx context.Context, username string) (bool, error)
}

type authenticatonServiceImpl struct {
	authenticatonRepo postgres.AuthenticatonRepository
	logger            *slog.Logger
}

func NewAuthenticatonService(authenticatonRepo postgres.AuthenticatonRepository, logger *slog.Logger) AuthService {
	return &authenticatonServiceImpl{
		authenticatonRepo: authenticatonRepo,
		logger: logger,
	}
}

func (s *authenticatonServiceImpl) Register(ctx context.Context, user models.Register) (*models.Response, error) {
	hashedPassword, err := token.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	resp, err := s.authenticatonRepo.Register(user)

	if err != nil {
		s.logger.Error("Error registering user", "error", err)
		return nil, err
	}

	return resp, nil
}

func (s *authenticatonServiceImpl) Login(ctx context.Context, login models.LoginRequest) (*models.User, error) {
	user, err := s.authenticatonRepo.Login(login)

	if err != nil {
		s.logger.Error("Error login user", "error", err)
		return nil, err
	}

	return user, nil
}

func (s *authenticatonServiceImpl) SaveRefreshToken(ctx context.Context, req models.TokenRequest) error {
	err := s.authenticatonRepo.SaveRefreshToken(req)

	if err != nil {
		s.logger.Error("Error saving refresh token", "error", err)
		return err
	}

	return nil
}

func (s *authenticatonServiceImpl) InvalidateRefreshToken(ctx context.Context, username string) error {
	err := s.authenticatonRepo.InvalidateRefreshToken(username)

	if err != nil {
		s.logger.Error("Error invalidating refresh token", "error", err)
		return err
	}

	return nil
}

func (s *authenticatonServiceImpl) IsRefreshTokenValid(ctx context.Context, username string) (bool, error) {
	valid, err := s.authenticatonRepo.IsRefreshTokenValid(username)

	if err != nil {
		s.logger.Error("Error checking refresh token validity", "error", err)
		return false, err
	}

	return valid, nil
}
