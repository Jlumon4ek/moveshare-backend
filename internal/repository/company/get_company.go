package company

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetCompany(ctx context.Context, userID int64) (*models.Company, error) {
	query := `
		SELECT 
			c.id, 
			c.user_id, 
			c.company_name, 
			u.email AS email_address, 
			c.address, 
			c.state, 
			c.mc_license_number,
			c.company_description, 
			c.contact_person, 
			c.phone_number, 
			c.city, 
			c.zip_code, 
			c.dot_number,
			c.created_at, 
			c.updated_at
		FROM companies c
		JOIN users u ON c.user_id = u.id
		WHERE c.user_id = $1
	`
	var company models.Company
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&company.ID,
		&company.UserID,
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
		return nil, err
	}
	return &company, nil
}
