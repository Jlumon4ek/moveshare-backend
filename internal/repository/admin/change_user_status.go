package admin

import (
	"context"
)

func (r *repository) ChangeUserStatus(ctx context.Context, userID int, newStatus string) error {
	query := `UPDATE users SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, newStatus, userID)
	if err != nil {
		return err
	}
	return nil
}
