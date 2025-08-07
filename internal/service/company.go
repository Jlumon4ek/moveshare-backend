package service

import (
	"context"
	"fmt"
	"moveshare/internal/dto"
	"moveshare/internal/models"
	"moveshare/internal/repository/company"
	"moveshare/internal/repository/user"
)

type CompanyService interface {
	GetCompany(ctx context.Context, userID int64) (*models.Company, error)
	UpdateCompany(ctx context.Context, userID int64, req dto.UpdateCompanyRequest) error
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
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if company == nil || company.ID == 0 {
		return &models.Company{
			EmailAddress: company.EmailAddress,
		}, nil
	}

	return company, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, userID int64, req dto.UpdateCompanyRequest) error {
	company := &models.Company{
		CompanyName:        req.CompanyName,
		Address:            req.Address,
		State:              req.State,
		MCLicenseNumber:    req.MCLicenseNumber,
		CompanyDescription: req.CompanyDescription,
		ContactPerson:      req.ContactPerson,
		PhoneNumber:        req.PhoneNumber,
		City:               req.City,
		ZipCode:            req.ZipCode,
		DotNumber:          req.DotNumber,
	}

	return s.companyRepo.UpdateCompany(ctx, userID, company)
}
