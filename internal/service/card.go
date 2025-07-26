package service

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"moveshare/internal/repository/card"
)

type CardService interface {
	CreateCard(ctx context.Context, userID int64, cardData *models.Card) error
	GetCard(ctx context.Context, userID, cardID int64) (*models.Card, error)
	GetUserCards(ctx context.Context, userID int64) ([]models.Card, error)
	SetDefaultCard(ctx context.Context, userID, cardID int64) error
	DeleteCard(ctx context.Context, userID, cardID int64) error
	UpdateCard(ctx context.Context, userID int64, cardData *models.Card) error
}

type cardService struct {
	cardRepo card.CardRepository
}

func NewCardService(cardRepo card.CardRepository) CardService {
	return &cardService{
		cardRepo: cardRepo,
	}
}

func (s *cardService) CreateCard(ctx context.Context, userID int64, cardData *models.Card) error {
	if cardData == nil {
		return fmt.Errorf("card data cannot be nil")
	}
	cardData.UserID = userID

	return s.cardRepo.CreateCard(ctx, cardData)
}

func (s *cardService) GetCard(ctx context.Context, userID, cardID int64) (*models.Card, error) {
	if userID <= 0 || cardID <= 0 {
		return nil, fmt.Errorf("invalid user ID or card ID")
	}

	return s.cardRepo.GetCardByID(ctx, userID, cardID)
}

func (s *cardService) GetUserCards(ctx context.Context, userID int64) ([]models.Card, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	return s.cardRepo.GetUserCards(ctx, userID)
}

func (s *cardService) SetDefaultCard(ctx context.Context, userID, cardID int64) error {
	if userID <= 0 || cardID <= 0 {
		return fmt.Errorf("invalid user ID or card ID")
	}

	card, err := s.cardRepo.GetCardByID(ctx, userID, cardID)
	if err != nil {
		return fmt.Errorf("failed to get card: %w", err)
	}
	if card == nil {
		return fmt.Errorf("card not found")
	}

	return s.cardRepo.SetDefaultCard(ctx, userID, cardID)
}

func (s *cardService) DeleteCard(ctx context.Context, userID, cardID int64) error {
	if userID <= 0 || cardID <= 0 {
		return fmt.Errorf("invalid user ID or card ID")
	}

	card, err := s.cardRepo.GetCardByID(ctx, userID, cardID)
	if err != nil {
		return fmt.Errorf("failed to get card: %w", err)
	}
	if card == nil {
		return fmt.Errorf("card not found")
	}

	return s.cardRepo.SoftDeleteCard(ctx, userID, cardID)
}

func (s *cardService) UpdateCard(ctx context.Context, userID int64, cardData *models.Card) error {
	if cardData == nil {
		return fmt.Errorf("card data cannot be nil")
	}
	if userID <= 0 || cardData.ID <= 0 {
		return fmt.Errorf("invalid user ID or card ID")
	}

	cardData.UserID = userID

	existingCard, err := s.cardRepo.GetCardByID(ctx, userID, cardData.ID)
	if err != nil {
		return fmt.Errorf("failed to get card: %w", err)
	}
	if existingCard == nil {
		return fmt.Errorf("card not found")
	}

	return s.cardRepo.UpdateCard(ctx, userID, cardData)
}
