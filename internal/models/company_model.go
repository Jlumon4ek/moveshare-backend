package models

import "time"

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
