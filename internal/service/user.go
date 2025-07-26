package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository/user"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
}

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) FindUserByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error) {
	user, err := s.userRepo.FindUserByEmailOrUsername(ctx, identifier)
	if err != nil {
		return nil, err
	}

	return user, nil
}
