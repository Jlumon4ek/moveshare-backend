package company

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepository interface {
	GetCompany(ctx context.Context, userID int64) (*models.Company, error)
	UpdateCompany(ctx context.Context, userID int64, company *models.Company) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) CompanyRepository {
	return &repository{db: db}
}
