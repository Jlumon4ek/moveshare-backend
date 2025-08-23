package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
	FindUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetUserByID(userID int64) (*models.User, error)
	UpdateProfilePhotoID(userID int64, photoID string) error
	CheckPassword(password, hash string) bool
	UpdatePassword(userID int64, newPassword string) error
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

func (s *userService) FindUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.userRepo.FindUserByID(ctx, userID)
}

func (s *userService) GetUserByID(userID int64) (*models.User, error) {
	ctx := context.Background()
	return s.userRepo.FindUserByID(ctx, userID)
}

func (s *userService) UpdateProfilePhotoID(userID int64, photoID string) error {
	ctx := context.Background()
	return s.userRepo.UpdateProfilePhotoID(ctx, userID, photoID)
}

func (s *userService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *userService) UpdatePassword(userID int64, newPassword string) error {
	ctx := context.Background()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}
