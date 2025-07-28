package dto

type CompanyResponse struct {
	CompanyName        string `json:"company_name" example:"Example Company"`
	EmailAddress       string `json:"email_address" example:"example@company.com"`
	Address            string `json:"address" example:"123 Main St"`
	State              string `json:"state" example:"CA"`
	MCLicenseNumber    string `json:"mc_license_number" example:"123456"`
	CompanyDescription string `json:"company_description" example:"A brief description of the company"`
	ContactPerson      string `json:"contact_person" example:"John Doe"`
	PhoneNumber        string `json:"phone_number" example:"(123) 456-7890"`
	City               string `json:"city" example:"Los Angeles"`
	ZipCode            string `json:"zip_code" example:"90001"`
	DotNumber          string `json:"dot_number" example:"DOT1234567"`
}
