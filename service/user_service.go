package service

import (
	"auth-service/api/token"
	pb "auth-service/generated/user"
	"auth-service/storage/postgres"
	"context"
	"log/slog"
)

type UserService interface {
	CreateUser(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error)
	UpdateUser(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	UpdateUserById(context.Context, *pb.UpdateUserByIdRequest) (*pb.UpdateUserByIdResponse, error)
	DeleteUser(context.Context, *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	GetAllUsers(context.Context, *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error)
}

type userService struct {
	pb.UnimplementedUserServiceServer
	userRepo postgres.UserRepository
	logger   *slog.Logger
}

func NewUserService(userRepo postgres.UserRepository, logger *slog.Logger) *userService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := token.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	in.Password = hashedPassword

	resp, err := s.userRepo.CreateUser(in)
	if err != nil {
		s.logger.Error("Error creating user", "error", err)
		return nil, err
	}

	return resp, err
}

func (s *userService) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.userRepo.GetUserByID(in.Id)
	if err != nil {
		s.logger.Error("Error getting user by ID", "error", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context, in *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.userRepo.GetAllUsers(in)
	if err != nil {
		s.logger.Error("Error getting all users", "error", err)
		return nil, err
	}

	return &pb.GetAllUsersResponse{Users: users}, nil
}

func (s *userService) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	resp, err := s.userRepo.UpdateUser(in)
	if err != nil {
		s.logger.Error("Error updating user", "error", err)
		return nil, err
	}

	return resp, nil
}

func (s *userService) UpdateUserById(ctx context.Context, in *pb.UpdateUserByIdRequest) (*pb.UpdateUserByIdResponse, error) {
	resp, err := s.userRepo.UpdateUserById(in)
	if err != nil {
		s.logger.Error("Error updating user by ID", "error", err)
		return nil, err
	}

	return resp, nil
}

func (s *userService) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	resp, err := s.userRepo.DeleteUser(in)
	if err != nil {
		s.logger.Error("Error deleting user", "error", err)
		return nil, err
	}

	return resp, nil
}

func (s *userService) DeleteUserByID(ctx context.Context, in *pb.DeleteByIdRequest) (*pb.DeleteUserResponse, error) {
	res, err := s.userRepo.DeleteUserByID(in)
	if err != nil {
		s.logger.Error("Error deleting user by ID", "error", err)
		return nil, err
	}
	return res, nil
}
