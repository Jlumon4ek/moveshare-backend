package password_reset

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetValidResetCode(ctx context.Context, email, code string) (*models.PasswordResetCode, error) {
	query := `
		SELECT id, user_id, email, code, expires_at, used, created_at, updated_at
		FROM password_reset_codes
		WHERE email = $1 AND code = $2 AND used = FALSE AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	resetCode := &models.PasswordResetCode{}
	err := r.db.QueryRow(ctx, query, email, code).Scan(
		&resetCode.ID,
		&resetCode.UserID,
		&resetCode.Email,
		&resetCode.Code,
		&resetCode.ExpiresAt,
		&resetCode.Used,
		&resetCode.CreatedAt,
		&resetCode.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return resetCode, nil
}