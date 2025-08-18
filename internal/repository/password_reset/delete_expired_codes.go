package password_reset

import (
	"context"
)

func (r *repository) DeleteExpiredCodes(ctx context.Context) error {
	query := `
		DELETE FROM password_reset_codes
		WHERE expires_at < NOW()
	`

	_, err := r.db.Exec(ctx, query)
	return err
}