package company

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetCompany(ctx context.Context, userID int64) (*models.Company, error) {
	query := `
		SELECT 
			COALESCE(c.id, 0) AS id,
			COALESCE(c.company_name, '') AS company_name, 
			u.email AS email_address,
			COALESCE(c.address, '') AS address, 
			COALESCE(c.state, '') AS state, 
			COALESCE(c.mc_license_number, '') AS mc_license_number,
			COALESCE(c.company_description, '') AS company_description, 
			COALESCE(c.contact_person, '') AS contact_person, 
			COALESCE(c.phone_number, '') AS phone_number, 
			COALESCE(c.city, '') AS city, 
			COALESCE(c.zip_code, '') AS zip_code, 
			COALESCE(c.dot_number, '') AS dot_number,
			COALESCE(c.created_at, now()) AS created_at, 
			COALESCE(c.updated_at, now()) AS updated_at
		FROM users u
		LEFT JOIN companies c ON c.user_id = u.id
		WHERE u.id = $1
	`

	var company models.Company
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&company.ID,
		&company.CompanyName,
		&company.EmailAddress,
		&company.Address,
		&company.State,
		&company.MCLicenseNumber,
		&company.CompanyDescription,
		&company.ContactPerson,
		&company.PhoneNumber,
		&company.City,
		&company.ZipCode,
		&company.DotNumber,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &company, nil
}
