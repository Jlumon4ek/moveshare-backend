package company

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) UpdateCompany(ctx context.Context, userID int64, company *models.Company) error {
	query := `
		INSERT INTO companies (
			user_id, company_name, address, state, mc_license_number,
			company_description, contact_person, phone_number, city, zip_code, dot_number
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id) DO UPDATE
		SET company_name = COALESCE($2, companies.company_name),
			address = COALESCE($3, companies.address),
			state = COALESCE($4, companies.state),
			mc_license_number = COALESCE($5, companies.mc_license_number),
			company_description = COALESCE($6, companies.company_description),
			contact_person = COALESCE($7, companies.contact_person),
			phone_number = COALESCE($8, companies.phone_number),
			city = COALESCE($9, companies.city),
			zip_code = COALESCE($10, companies.zip_code),
			dot_number = COALESCE($11, companies.dot_number),
			updated_at = NOW()
	`
	_, err := r.db.Exec(ctx, query,
		userID, company.CompanyName, company.Address, company.State,
		company.MCLicenseNumber, company.CompanyDescription, company.ContactPerson,
		company.PhoneNumber, company.City, company.ZipCode, company.DotNumber,
	)
	return err
}
