package email_verification

import "context"

func (r *repository) MarkCodeAsUsed(ctx context.Context, id int64) error {
	query := `
		UPDATE email_verifications
		SET used = true
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query, id)
	return err
}