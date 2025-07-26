package verification

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationRepository interface {
	InsertFileID(ctx context.Context, userID int64, objectName string, fileType string) error
	SelectVerificationFiles(ctx context.Context, userID int64) ([]models.VerificationFile, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewVerificationRepository(db *pgxpool.Pool) VerificationRepository {
	return &repository{db: db}
}
