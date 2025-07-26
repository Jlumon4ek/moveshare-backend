package models

import "time"

type Card struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CardNumber  string    `json:"card_number"`
	CardHolder  string    `json:"card_holder"`
	ExpiryMonth int       `json:"expiry_month"`
	ExpiryYear  int       `json:"expiry_year"`
	CVV         string    `json:"cvv,omitempty"`
	CardType    string    `json:"card_type"`
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
