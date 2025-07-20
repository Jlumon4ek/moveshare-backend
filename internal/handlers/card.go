package handlers

import (
	"encoding/json"
	"moveshare/internal/repository"
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// CardHandler handles HTTP requests for card operations
type CardHandler struct {
	cardService service.CardService
}

// NewCardHandler creates a new CardHandler
func NewCardHandler(cardService service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

// CreateCardRequest represents the card creation request payload
type CreateCardRequest struct {
	CardNumber  string `json:"card_number" example:"4111111111111111"`
	CardHolder  string `json:"card_holder" example:"John Doe"`
	ExpiryMonth int    `json:"expiry_month" example:"12"`
	ExpiryYear  int    `json:"expiry_year" example:"2025"`
	CVV         string `json:"cvv" example:"123"`
	IsDefault   bool   `json:"is_default" example:"false"`
}

// UpdateCardRequest represents the card update request payload
type UpdateCardRequest struct {
	CardHolder  *string `json:"card_holder,omitempty" example:"Jane Doe"`
	ExpiryMonth *int    `json:"expiry_month,omitempty" example:"11"`
	ExpiryYear  *int    `json:"expiry_year,omitempty" example:"2026"`
	IsDefault   *bool   `json:"is_default,omitempty" example:"true"`
}

// CardResponse represents a single card response
type CardResponse struct {
	Card *repository.Card `json:"card"`
}

// CardsResponse represents multiple cards response
type CardsResponse struct {
	Cards []repository.Card `json:"cards"`
}

// CreateCard handles card creation
// @Summary Create a new payment card
// @Description Creates a new payment card for the current user
// @Tags Payment Settings
// @Accept json
// @Produce json
// @Param body body CreateCardRequest true "Card data"
// @Success 201 {object} CardResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cards [post]
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	card := &repository.Card{
		CardNumber:  req.CardNumber,
		CardHolder:  req.CardHolder,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		CVV:         req.CVV,
		IsDefault:   req.IsDefault,
	}

	err := h.cardService.CreateCard(r.Context(), userID, card)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Clear CVV from response for security
	card.CVV = ""

	resp := CardResponse{Card: card}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetUserCards handles fetching all user cards
// @Summary Get user payment cards
// @Description Fetches all active payment cards for the current user
// @Tags Payment Settings
// @Produce json
// @Success 200 {object} CardsResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cards [get]
func (h *CardHandler) GetUserCards(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	cards, err := h.cardService.GetUserCards(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := CardsResponse{Cards: cards}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCardByID handles fetching a specific card
// @Summary Get a specific payment card
// @Description Fetches a specific payment card by ID for the current user
// @Tags Payment Settings
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} CardResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cards/{id} [get]
func (h *CardHandler) GetCardByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	cardIDStr := chi.URLParam(r, "id")
	if cardIDStr == "" {
		http.Error(w, `{"error": "card ID is required"}`, http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseInt(cardIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid card ID"}`, http.StatusBadRequest)
		return
	}

	card, err := h.cardService.GetCardByID(r.Context(), userID, cardID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if card == nil {
		http.Error(w, `{"error": "card not found"}`, http.StatusNotFound)
		return
	}

	resp := CardResponse{Card: card}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateCard handles card updates
// @Summary Update a payment card
// @Description Updates a payment card for the current user
// @Tags Payment Settings
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param body body UpdateCardRequest true "Card update data"
// @Success 200 {object} CardResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cards/{id} [put]
func (h *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	cardIDStr := chi.URLParam(r, "id")
	if cardIDStr == "" {
		http.Error(w, `{"error": "card ID is required"}`, http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseInt(cardIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid card ID"}`, http.StatusBadRequest)
		return
	}

	var req UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	card := &repository.Card{ID: cardID}
	if req.CardHolder != nil {
		card.CardHolder = *req.CardHolder
	}
	if req.ExpiryMonth != nil {
		card.ExpiryMonth = *req.ExpiryMonth
	}
	if req.ExpiryYear != nil {
		card.ExpiryYear = *req.ExpiryYear
	}
	if req.IsDefault != nil {
		card.IsDefault = *req.IsDefault
	}

	err = h.cardService.UpdateCard(r.Context(), userID, card)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Fetch updated card
	updatedCard, err := h.cardService.GetCardByID(r.Context(), userID, cardID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := CardResponse{Card: updatedCard}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteCard handles card soft deletion
// @Summary Delete a payment card
// @Description Soft deletes a payment card for the current user
// @Tags Payment Settings
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cards/{id} [delete]
func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	cardIDStr := chi.URLParam(r, "id")
	if cardIDStr == "" {
		http.Error(w, `{"error": "card ID is required"}`, http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseInt(cardIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid card ID"}`, http.StatusBadRequest)
		return
	}

	err = h.cardService.DeleteCard(r.Context(), userID, cardID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Card deleted successfully"})
}

// SetDefaultCard handles setting a card as default
// @Summary Set card as default
// @Description Sets a payment card as the default payment method
// @Tags Payment Settings
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cards/{id}/default [post]
func (h *CardHandler) SetDefaultCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	cardIDStr := chi.URLParam(r, "id")
	if cardIDStr == "" {
		http.Error(w, `{"error": "card ID is required"}`, http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseInt(cardIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid card ID"}`, http.StatusBadRequest)
		return
	}

	err = h.cardService.SetDefaultCard(r.Context(), userID, cardID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Default card updated successfully"})
}
