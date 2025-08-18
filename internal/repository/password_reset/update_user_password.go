package password_reset

import (
	"context"
)

func (r *repository) UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, hashedPassword, userID)
	return err
}