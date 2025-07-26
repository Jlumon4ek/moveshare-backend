package card

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *repository) UpdateCard(ctx context.Context, userID int64, card *models.Card) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if card.IsDefault {
		_, err = tx.Exec(ctx,
			"UPDATE cards SET is_default = false WHERE user_id = $1 AND is_active = true AND id != $2",
			userID, card.ID)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE cards 
		SET card_holder = COALESCE($3, card_holder),
		    expiry_month = COALESCE($4, expiry_month),
		    expiry_year = COALESCE($5, expiry_year),
		    card_type = COALESCE($6, card_type),
		    is_default = COALESCE($7, is_default),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2 AND is_active = true
		RETURNING updated_at
	`
	err = tx.QueryRow(ctx, query,
		card.ID, userID, card.CardHolder, card.ExpiryMonth, card.ExpiryYear,
		card.CardType, card.IsDefault,
	).Scan(&card.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}
