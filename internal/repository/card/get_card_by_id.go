package card

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetCardByID(ctx context.Context, userID, cardID int64) (*models.Card, error) {
	query := `
		SELECT id, user_id, card_number, card_holder, expiry_month, expiry_year,
		       card_type, is_default, is_active, created_at, updated_at
		FROM cards
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`
	var card models.Card
	err := r.db.QueryRow(ctx, query, cardID, userID).Scan(
		&card.ID, &card.UserID, &card.CardNumber, &card.CardHolder,
		&card.ExpiryMonth, &card.ExpiryYear, &card.CardType,
		&card.IsDefault, &card.IsActive, &card.CreatedAt, &card.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &card, nil
}
