package schemas

import "moveshare/internal/models"

type CompanyUpdateRequest struct {
	CompanyName        *string `json:"company_name"`
	EmailAddress       *string `json:"email_address"`
	Address            *string `json:"address"`
	State              *string `json:"state"`
	MCLicenseNumber    *string `json:"mc_license_number"`
	CompanyDescription *string `json:"company_description"`
	ContactPerson      *string `json:"contact_person"`
	PhoneNumber        *string `json:"phone_number"`
	City               *string `json:"city"`
	ZipCode            *string `json:"zip_code"`
	DotNumber          *string `json:"dot_number"`
}

type CompanyResponse struct {
	Company *models.Company `json:"company"`
}
