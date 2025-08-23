package email_verification

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) CreateVerificationCode(ctx context.Context, verification *models.EmailVerification) error {
	query := `
		INSERT INTO email_verifications (email, code, expires_at, created_at, used)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	
	err := r.db.QueryRow(ctx, query,
		verification.Email,
		verification.Code,
		verification.ExpiresAt,
		verification.CreatedAt,
		verification.Used,
	).Scan(&verification.ID)
	
	return err
}