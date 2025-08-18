package password_reset

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) CreateResetCode(ctx context.Context, code *models.PasswordResetCode) error {
	query := `
		INSERT INTO password_reset_codes (user_id, email, code, expires_at, used)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		code.UserID,
		code.Email,
		code.Code,
		code.ExpiresAt,
		code.Used,
	).Scan(&code.ID, &code.CreatedAt, &code.UpdatedAt)
}