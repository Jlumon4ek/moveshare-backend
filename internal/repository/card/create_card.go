package card

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) CreateCard(ctx context.Context, card *models.Card) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if card.IsDefault {
		_, err = tx.Exec(ctx,
			"UPDATE cards SET is_default = false WHERE user_id = $1 AND is_active = true",
			card.UserID)
		if err != nil {
			return err
		}
	} else {
		var count int
		err = tx.QueryRow(ctx,
			"SELECT COUNT(*) FROM cards WHERE user_id = $1 AND is_active = true",
			card.UserID).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			card.IsDefault = true
		}
	}

	query := `
		INSERT INTO cards (user_id, card_number, card_holder, expiry_month, expiry_year, 
		                  cvv, card_type, is_default, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		card.UserID, card.CardNumber, card.CardHolder, card.ExpiryMonth, card.ExpiryYear,
		card.CVV, card.CardType, card.IsDefault,
	).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
