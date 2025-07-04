package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"moveshare/internal/auth"
	"moveshare/internal/repository"
)

// UserService defines the interface for user business logic
type UserService interface {
	SignUp(ctx context.Context, username, email, password string) error
	SignIn(ctx context.Context, email, password string) (*repository.User, string, string, error)
}

// userService implements UserService
type userService struct {
	userRepo repository.UserRepository
	jwtAuth  auth.JWTAuth
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository, jwtAuth auth.JWTAuth) UserService {
	return &userService{userRepo: userRepo, jwtAuth: jwtAuth}
}

// SignUp creates a new user
func (s *userService) SignUp(ctx context.Context, username, email, password string) error {
	user := &repository.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	return s.userRepo.CreateUser(ctx, user)
}

// SignIn authenticates a user and returns access and refresh tokens
func (s *userService) SignIn(ctx context.Context, email, password string) (*repository.User, string, string, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", "", errors.New("invalid password")
	}

	accessToken, err := s.jwtAuth.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.jwtAuth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
