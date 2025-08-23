package email_verification

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailVerificationRepository interface {
	CreateVerificationCode(ctx context.Context, verification *models.EmailVerification) error
	GetVerificationByEmailAndCode(ctx context.Context, email, code string) (*models.EmailVerification, error)
	MarkCodeAsUsed(ctx context.Context, id int64) error
	DeleteExpiredCodes(ctx context.Context) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewEmailVerificationRepository(db *pgxpool.Pool) EmailVerificationRepository {
	return &repository{db: db}
}