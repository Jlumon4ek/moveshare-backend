package card

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepository interface {
	CreateCard(ctx context.Context, card *models.Card) error
	GetCardByID(ctx context.Context, userID, cardID int64) (*models.Card, error)
	GetUserCards(ctx context.Context, userID int64) ([]models.Card, error)
	SetDefaultCard(ctx context.Context, userID, cardID int64) error
	SoftDeleteCard(ctx context.Context, userID, cardID int64) error
	UpdateCard(ctx context.Context, userID int64, card *models.Card) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) CardRepository {
	return &repository{db: db}
}
