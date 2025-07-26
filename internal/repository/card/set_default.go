package card

import "context"

func (r *repository) SetDefaultCard(ctx context.Context, userID, cardID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"UPDATE cards SET is_default = false WHERE user_id = $1 AND is_active = true",
		userID)
	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx,
		"UPDATE cards SET is_default = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2 AND is_active = true",
		cardID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return err
	}

	return tx.Commit(ctx)
}
