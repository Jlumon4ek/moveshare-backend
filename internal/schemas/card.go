package schemas

import "moveshare/internal/models"

type CreateCardRequest struct {
	CardNumber  string `json:"card_number"`
	CardHolder  string `json:"card_holder"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	IsDefault   bool   `json:"is_default"`
}

type UpdateCardRequest struct {
	CardHolder  *string `json:"card_holder,omitempty"`
	ExpiryMonth *int    `json:"expiry_month,omitempty"`
	ExpiryYear  *int    `json:"expiry_year,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
}

type CardResponse struct {
	Card *models.Card `json:"card"`
}

type CardsResponse struct {
	Cards []models.Card `json:"cards"`
}
