package email_verification

import "context"

func (r *repository) DeleteExpiredCodes(ctx context.Context) error {
	query := `
		DELETE FROM email_verifications
		WHERE expires_at < NOW()
	`
	
	_, err := r.db.Exec(ctx, query)
	return err
}