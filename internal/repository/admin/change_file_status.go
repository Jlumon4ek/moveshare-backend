package admin

import "context"

func (r *repository) ChangeVerificationFileStatus(ctx context.Context, fileID int, newStatus string) error {
	query := `UPDATE verification_file SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, newStatus, fileID)
	if err != nil {
		return err
	}
	return nil
}
