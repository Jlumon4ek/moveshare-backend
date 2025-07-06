package service

import (
	"context"
	"errors"
	"moveshare/internal/repository"
)

// CompanyService defines the interface for company business logic
type CompanyService interface {
	GetCompany(ctx context.Context, userID int64) (*repository.Company, error)
	UpdateCompany(ctx context.Context, userID int64, company *repository.Company) error
}

// companyService implements CompanyService
type companyService struct {
	companyRepo repository.CompanyRepository
}

// NewCompanyService creates a new CompanyService
func NewCompanyService(companyRepo repository.CompanyRepository) CompanyService {
	return &companyService{companyRepo: companyRepo}
}

// GetCompany fetches the company data for the given userID
func (s *companyService) GetCompany(ctx context.Context, userID int64) (*repository.Company, error) {
	return s.companyRepo.GetCompany(ctx, userID)
}

// UpdateCompany updates the company data for the given userID
func (s *companyService) UpdateCompany(ctx context.Context, userID int64, company *repository.Company) error {
	if company == nil {
		return errors.New("company data is required")
	}
	company.UserID = userID // Ensure the userID is set
	return s.companyRepo.UpdateCompany(ctx, userID, company)
}
