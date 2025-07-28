package service

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"moveshare/internal/repository/company"
	"moveshare/internal/repository/user"
)

type CompanyService interface {
	GetCompany(ctx context.Context, userID int64) (*models.Company, error)
	UpdateCompany(ctx context.Context, userID int64, company *models.Company) error
}

type companyService struct {
	companyRepo company.CompanyRepository
	userRepo    user.UserRepository
}

func NewCompanyService(companyRepo company.CompanyRepository, userRepo user.UserRepository) CompanyService {
	return &companyService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

func (s *companyService) GetCompany(ctx context.Context, userID int64) (*models.Company, error) {
	company, err := s.companyRepo.GetCompany(ctx, userID)
	if err != nil {
		return nil, err
	}

	if company == nil {
		user, err := s.userRepo.FindUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user: %w", err)
		}
		company = &models.Company{
			EmailAddress: user.Email,
		}
	}

	return company, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, userID int64, company *models.Company) error {
	return s.companyRepo.UpdateCompany(ctx, userID, company)
}
