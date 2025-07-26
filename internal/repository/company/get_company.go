package company

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetCompany(ctx context.Context, userID int64) (*models.Company, error) {
	query := `
		SELECT id, user_id, company_name, email_address, address, state, mc_license_number,
		       company_description, contact_person, phone_number, city, zip_code, dot_number,
		       created_at, updated_at
		FROM companies
		WHERE user_id = $1
	`
	var company models.Company
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&company.ID, &company.UserID, &company.CompanyName, &company.EmailAddress,
		&company.Address, &company.State, &company.MCLicenseNumber, &company.CompanyDescription,
		&company.ContactPerson, &company.PhoneNumber, &company.City, &company.ZipCode,
		&company.DotNumber, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}
