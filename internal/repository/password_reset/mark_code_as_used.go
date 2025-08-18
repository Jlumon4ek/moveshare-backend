package password_reset

import (
	"context"
)

func (r *repository) MarkCodeAsUsed(ctx context.Context, id int) error {
	query := `
		UPDATE password_reset_codes
		SET used = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}