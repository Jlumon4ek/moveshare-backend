package card

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *repository) SoftDeleteCard(ctx context.Context, userID, cardID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var isDefault bool
	err = tx.QueryRow(ctx,
		"SELECT is_default FROM cards WHERE id = $1 AND user_id = $2 AND is_active = true",
		cardID, userID).Scan(&isDefault)
	if err != nil {
		if err == pgx.ErrNoRows {
			return err
		}
		return err
	}

	result, err := tx.Exec(ctx,
		"UPDATE cards SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2",
		cardID, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return err
	}

	if isDefault {
		_, err = tx.Exec(ctx, `
			UPDATE cards 
			SET is_default = true 
			WHERE user_id = $1 AND is_active = true 
			AND id = (
				SELECT id FROM cards 
				WHERE user_id = $1 AND is_active = true 
				ORDER BY created_at ASC 
				LIMIT 1
			)`,
			userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
