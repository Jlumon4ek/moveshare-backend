package service

import (
	"context"
	"errors"
	"moveshare/internal/repository"
	"regexp"
	"strconv"
	"time"
)

// CardService defines the interface for card business logic
type CardService interface {
	CreateCard(ctx context.Context, userID int64, card *repository.Card) error
	GetUserCards(ctx context.Context, userID int64) ([]repository.Card, error)
	GetCardByID(ctx context.Context, userID, cardID int64) (*repository.Card, error)
	UpdateCard(ctx context.Context, userID int64, card *repository.Card) error
	DeleteCard(ctx context.Context, userID, cardID int64) error
	SetDefaultCard(ctx context.Context, userID, cardID int64) error
}

// cardService implements CardService
type cardService struct {
	cardRepo repository.CardRepository
}

// NewCardService creates a new CardService
func NewCardService(cardRepo repository.CardRepository) CardService {
	return &cardService{cardRepo: cardRepo}
}

// CreateCard creates a new card with validation
func (s *cardService) CreateCard(ctx context.Context, userID int64, card *repository.Card) error {
	if card == nil {
		return errors.New("card data is required")
	}

	// Validate required fields
	if err := s.validateCard(card); err != nil {
		return err
	}

	// Detect card type
	card.CardType = detectCardType(card.CardNumber)
	card.UserID = userID

	return s.cardRepo.CreateCard(ctx, card)
}

// GetUserCards retrieves all active cards for a user
func (s *cardService) GetUserCards(ctx context.Context, userID int64) ([]repository.Card, error) {
	return s.cardRepo.GetUserCards(ctx, userID)
}

// GetCardByID retrieves a specific card by ID for a user
func (s *cardService) GetCardByID(ctx context.Context, userID, cardID int64) (*repository.Card, error) {
	return s.cardRepo.GetCardByID(ctx, userID, cardID)
}

// UpdateCard updates an existing card
func (s *cardService) UpdateCard(ctx context.Context, userID int64, card *repository.Card) error {
	if card == nil {
		return errors.New("card data is required")
	}

	// Validate expiry date if provided
	if card.ExpiryMonth != 0 && card.ExpiryYear != 0 {
		if err := s.validateExpiryDate(card.ExpiryMonth, card.ExpiryYear); err != nil {
			return err
		}
	}

	return s.cardRepo.UpdateCard(ctx, userID, card)
}

// DeleteCard performs soft delete on a card
func (s *cardService) DeleteCard(ctx context.Context, userID, cardID int64) error {
	return s.cardRepo.SoftDeleteCard(ctx, userID, cardID)
}

// SetDefaultCard sets a card as the default payment method
func (s *cardService) SetDefaultCard(ctx context.Context, userID, cardID int64) error {
	return s.cardRepo.SetDefaultCard(ctx, userID, cardID)
}

// validateCard validates card data
func (s *cardService) validateCard(card *repository.Card) error {
	// Validate card number
	if card.CardNumber == "" {
		return errors.New("card number is required")
	}

	// Remove spaces and validate format
	cardNumber := regexp.MustCompile(`\s+`).ReplaceAllString(card.CardNumber, "")
	if !isValidCardNumber(cardNumber) {
		return errors.New("invalid card number")
	}
	card.CardNumber = cardNumber

	// Validate card holder
	if card.CardHolder == "" {
		return errors.New("card holder name is required")
	}

	// Validate expiry date
	if err := s.validateExpiryDate(card.ExpiryMonth, card.ExpiryYear); err != nil {
		return err
	}

	// Validate CVV
	if card.CVV == "" {
		return errors.New("CVV is required")
	}
	if len(card.CVV) < 3 || len(card.CVV) > 4 {
		return errors.New("CVV must be 3 or 4 digits")
	}

	return nil
}

// validateExpiryDate validates expiry month and year
func (s *cardService) validateExpiryDate(month, year int) error {
	if month < 1 || month > 12 {
		return errors.New("expiry month must be between 1 and 12")
	}

	currentYear := time.Now().Year()
	if year < currentYear || year > currentYear+20 {
		return errors.New("invalid expiry year")
	}

	// Check if card is expired
	currentMonth := int(time.Now().Month())
	if year == currentYear && month < currentMonth {
		return errors.New("card has expired")
	}

	return nil
}

// isValidCardNumber validates card number using Luhn algorithm
func isValidCardNumber(cardNumber string) bool {
	// Check if all characters are digits
	if !regexp.MustCompile(`^\d+$`).MatchString(cardNumber) {
		return false
	}

	// Check length
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	// Luhn algorithm
	sum := 0
	alternate := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cardNumber[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// detectCardType detects card type based on card number
func detectCardType(cardNumber string) string {
	// Visa
	if regexp.MustCompile(`^4`).MatchString(cardNumber) {
		return "Visa"
	}
	// MasterCard
	if regexp.MustCompile(`^5[1-5]`).MatchString(cardNumber) {
		return "MasterCard"
	}
	// American Express
	if regexp.MustCompile(`^3[47]`).MatchString(cardNumber) {
		return "American Express"
	}
	// Discover
	if regexp.MustCompile(`^6(?:011|5)`).MatchString(cardNumber) {
		return "Discover"
	}

	return "Unknown"
}
