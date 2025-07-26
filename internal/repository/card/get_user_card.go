package card

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUserCards(ctx context.Context, userID int64) ([]models.Card, error) {
	query := `
		SELECT id, user_id, card_number, card_holder, expiry_month, expiry_year,
		       card_type, is_default, is_active, created_at, updated_at
		FROM cards
		WHERE user_id = $1 AND is_active = true
		ORDER BY is_default DESC, created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		err := rows.Scan(
			&card.ID, &card.UserID, &card.CardNumber, &card.CardHolder,
			&card.ExpiryMonth, &card.ExpiryYear, &card.CardType,
			&card.IsDefault, &card.IsActive, &card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, rows.Err()
}
