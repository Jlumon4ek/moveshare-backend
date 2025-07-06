package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Company represents a company entity
type Company struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	CompanyName        string    `json:"company_name"`
	EmailAddress       string    `json:"email_address"`
	Address            string    `json:"address"`
	State              string    `json:"state"`
	MCLicenseNumber    string    `json:"mc_license_number"`
	CompanyDescription string    `json:"company_description"`
	ContactPerson      string    `json:"contact_person"`
	PhoneNumber        string    `json:"phone_number"`
	City               string    `json:"city"`
	ZipCode            string    `json:"zip_code"`
	DotNumber          string    `json:"dot_number"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CompanyRepository defines the interface for company data operations
type CompanyRepository interface {
	GetCompany(ctx context.Context, userID int64) (*Company, error)
	UpdateCompany(ctx context.Context, userID int64, company *Company) error
}

// companyRepository implements CompanyRepository
type companyRepository struct {
	db *pgxpool.Pool
}

// NewCompanyRepository creates a new CompanyRepository
func NewCompanyRepository(db *pgxpool.Pool) CompanyRepository {
	return &companyRepository{db: db}
}

// GetCompany fetches the company data for the given userID
func (r *companyRepository) GetCompany(ctx context.Context, userID int64) (*Company, error) {
	query := `
		SELECT id, user_id, company_name, email_address, address, state, mc_license_number,
		       company_description, contact_person, phone_number, city, zip_code, dot_number,
		       created_at, updated_at
		FROM companies
		WHERE user_id = $1
	`
	var company Company
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&company.ID, &company.UserID, &company.CompanyName, &company.EmailAddress,
		&company.Address, &company.State, &company.MCLicenseNumber, &company.CompanyDescription,
		&company.ContactPerson, &company.PhoneNumber, &company.City, &company.ZipCode,
		&company.DotNumber, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil // Return nil if no company exists
		}
		return nil, err
	}
	return &company, nil
}

// UpdateCompany updates the company data for the given userID
func (r *companyRepository) UpdateCompany(ctx context.Context, userID int64, company *Company) error {
	query := `
		INSERT INTO companies (user_id, company_name, email_address, address, state, mc_license_number,
		                      company_description, contact_person, phone_number, city, zip_code, dot_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (user_id) DO UPDATE
		SET company_name = COALESCE($2, companies.company_name),
		    email_address = COALESCE($3, companies.email_address),
		    address = COALESCE($4, companies.address),
		    state = COALESCE($5, companies.state),
		    mc_license_number = COALESCE($6, companies.mc_license_number),
		    company_description = COALESCE($7, companies.company_description),
		    contact_person = COALESCE($8, companies.contact_person),
		    phone_number = COALESCE($9, companies.phone_number),
		    city = COALESCE($10, companies.city),
		    zip_code = COALESCE($11, companies.zip_code),
		    dot_number = COALESCE($12, companies.dot_number),
		    updated_at = NOW()
	`
	_, err := r.db.Exec(ctx, query,
		userID, company.CompanyName, company.EmailAddress, company.Address, company.State,
		company.MCLicenseNumber, company.CompanyDescription, company.ContactPerson,
		company.PhoneNumber, company.City, company.ZipCode, company.DotNumber,
	)
	return err
}
