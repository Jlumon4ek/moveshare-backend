package user

import "context"

func (r *repository) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, newPassword, userID)
	return err
}
