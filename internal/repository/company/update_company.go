package company

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) UpdateCompany(ctx context.Context, userID int64, c *models.Company) error {
	query := `
		INSERT INTO companies (
			user_id, company_name, address, state, mc_license_number, 
			company_description, contact_person, phone_number, city, zip_code, dot_number, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, now(), now()
		)
		ON CONFLICT (user_id) DO UPDATE SET
			company_name = EXCLUDED.company_name,
			address = EXCLUDED.address,
			state = EXCLUDED.state,
			mc_license_number = EXCLUDED.mc_license_number,
			company_description = EXCLUDED.company_description,
			contact_person = EXCLUDED.contact_person,
			phone_number = EXCLUDED.phone_number,
			city = EXCLUDED.city,
			zip_code = EXCLUDED.zip_code,
			dot_number = EXCLUDED.dot_number,
			updated_at = now()
	`

	_, err := r.db.Exec(ctx, query,
		userID,
		c.CompanyName,
		c.Address,
		c.State,
		c.MCLicenseNumber,
		c.CompanyDescription,
		c.ContactPerson,
		c.PhoneNumber,
		c.City,
		c.ZipCode,
		c.DotNumber,
	)

	return err
}
