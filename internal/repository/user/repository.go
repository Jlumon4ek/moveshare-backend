package user

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
	FindUserByID(ctx context.Context, userID int64) (*models.User, error)
	UpdateProfilePhotoID(ctx context.Context, userID int64, photoID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &repository{db: db}
}
