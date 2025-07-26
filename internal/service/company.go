package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository/company"
)

type CompanyService interface {
	GetCompany(ctx context.Context, userID int64) (*models.Company, error)
	UpdateCompany(ctx context.Context, userID int64, company *models.Company) error
}

type companyService struct {
	companyRepo company.CompanyRepository
}

func NewCompanyService(companyRepo company.CompanyRepository) CompanyService {
	return &companyService{
		companyRepo: companyRepo,
	}
}

func (s *companyService) GetCompany(ctx context.Context, userID int64) (*models.Company, error) {
	return s.companyRepo.GetCompany(ctx, userID)
}

func (s *companyService) UpdateCompany(ctx context.Context, userID int64, company *models.Company) error {
	return s.companyRepo.UpdateCompany(ctx, userID, company)
}
