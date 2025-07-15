package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Card represents a card entity
type Card struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CardNumber  string    `json:"card_number"`
	CardHolder  string    `json:"card_holder"`
	ExpiryMonth int       `json:"expiry_month"`
	ExpiryYear  int       `json:"expiry_year"`
	CVV         string    `json:"cvv,omitempty"` // Не возвращаем в JSON для безопасности
	CardType    string    `json:"card_type"`     // Visa, MasterCard, etc.
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"` // Soft delete flag
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CardRepository defines the interface for card data operations
type CardRepository interface {
	CreateCard(ctx context.Context, card *Card) error
	GetUserCards(ctx context.Context, userID int64) ([]Card, error)
	GetCardByID(ctx context.Context, userID, cardID int64) (*Card, error)
	UpdateCard(ctx context.Context, userID int64, card *Card) error
	SoftDeleteCard(ctx context.Context, userID, cardID int64) error
	SetDefaultCard(ctx context.Context, userID, cardID int64) error
}

// cardRepository implements CardRepository
type cardRepository struct {
	db *pgxpool.Pool
}

// NewCardRepository creates a new CardRepository
func NewCardRepository(db *pgxpool.Pool) CardRepository {
	return &cardRepository{db: db}
}

// CreateCard creates a new card
func (r *cardRepository) CreateCard(ctx context.Context, card *Card) error {
	// Begin transaction to handle default card logic
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

func (r *cardRepository) GetUserCards(ctx context.Context, userID int64) ([]Card, error) {
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

	var cards []Card
	for rows.Next() {
		var card Card
		err := rows.Scan(
			&card.ID, &card.UserID, &card.CardNumber, &card.CardHolder,
			&card.ExpiryMonth, &card.ExpiryYear, &card.CardType,
			&card.IsDefault, &card.IsActive, &card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		card.CardNumber = maskCardNumber(card.CardNumber)
		cards = append(cards, card)
	}

	return cards, rows.Err()
}

func (r *cardRepository) GetCardByID(ctx context.Context, userID, cardID int64) (*Card, error) {
	query := `
		SELECT id, user_id, card_number, card_holder, expiry_month, expiry_year,
		       card_type, is_default, is_active, created_at, updated_at
		FROM cards
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`
	var card Card
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

	card.CardNumber = maskCardNumber(card.CardNumber)
	return &card, nil
}

// UpdateCard updates an existing card
func (r *cardRepository) UpdateCard(ctx context.Context, userID int64, card *Card) error {
	// Begin transaction to handle default card logic
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// If setting as default, remove default from other cards
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

// SoftDeleteCard performs soft delete on a card
func (r *cardRepository) SoftDeleteCard(ctx context.Context, userID, cardID int64) error {
	// Begin transaction to handle default card reassignment
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if the card to delete is default
	var isDefault bool
	err = tx.QueryRow(ctx,
		"SELECT is_default FROM cards WHERE id = $1 AND user_id = $2 AND is_active = true",
		cardID, userID).Scan(&isDefault)
	if err != nil {
		if err == pgx.ErrNoRows {
			return err // Card not found or already deleted
		}
		return err
	}

	// Soft delete the card
	result, err := tx.Exec(ctx,
		"UPDATE cards SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2",
		cardID, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return err // No rows updated
	}

	// If deleted card was default, set another card as default
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

// SetDefaultCard sets a card as the default payment method
func (r *cardRepository) SetDefaultCard(ctx context.Context, userID, cardID int64) error {
	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Remove default from all cards
	_, err = tx.Exec(ctx,
		"UPDATE cards SET is_default = false WHERE user_id = $1 AND is_active = true",
		userID)
	if err != nil {
		return err
	}

	// Set new default card
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

func maskCardNumber(cardNumber string) string {
	if len(cardNumber) <= 4 {
		return cardNumber
	}
	masked := ""
	for i := 0; i < len(cardNumber)-4; i++ {
		masked += "*"
	}
	return masked + cardNumber[len(cardNumber)-4:]
}
