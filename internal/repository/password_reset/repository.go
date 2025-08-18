package password_reset

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PasswordResetRepository interface {
	CreateResetCode(ctx context.Context, code *models.PasswordResetCode) error
	GetValidResetCode(ctx context.Context, email, code string) (*models.PasswordResetCode, error)
	MarkCodeAsUsed(ctx context.Context, id int) error
	DeleteExpiredCodes(ctx context.Context) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewPasswordResetRepository(db *pgxpool.Pool) PasswordResetRepository {
	return &repository{db: db}
}