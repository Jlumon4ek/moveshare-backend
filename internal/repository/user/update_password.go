package user

import "context"

func (r *repository) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`
	var updatedID int64
	err := r.db.QueryRow(ctx, query, newPassword, userID).Scan(&updatedID)
	if err != nil {
		return err
	}
	return nil
}
