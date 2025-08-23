package email_verification

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetVerificationByEmailAndCode(ctx context.Context, email, code string) (*models.EmailVerification, error) {
	query := `
		SELECT id, email, code, expires_at, created_at, used
		FROM email_verifications
		WHERE email = $1 AND code = $2 AND used = false AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	var verification models.EmailVerification
	err := r.db.QueryRow(ctx, query, email, code).Scan(
		&verification.ID,
		&verification.Email,
		&verification.Code,
		&verification.ExpiresAt,
		&verification.CreatedAt,
		&verification.Used,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &verification, nil
}