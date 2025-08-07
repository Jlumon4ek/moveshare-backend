package dto

import "moveshare/internal/models"

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

func NewCompanyResponse(c *models.Company) CompanyResponse {
	return CompanyResponse{
		CompanyName:        c.CompanyName,
		EmailAddress:       c.EmailAddress,
		Address:            c.Address,
		State:              c.State,
		MCLicenseNumber:    c.MCLicenseNumber,
		CompanyDescription: c.CompanyDescription,
		ContactPerson:      c.ContactPerson,
		PhoneNumber:        c.PhoneNumber,
		City:               c.City,
		ZipCode:            c.ZipCode,
		DotNumber:          c.DotNumber,
	}
}

type UpdateCompanyRequest struct {
	CompanyName        string `json:"company_name" example:"Example Company"`
	Address            string `json:"address" example:"123 Main St"`
	State              string `json:"state" example:"CA"`
	MCLicenseNumber    string `json:"mc_license_number" example:"123456"`
	CompanyDescription string `json:"company_description" example:"A short description"`
	ContactPerson      string `json:"contact_person" example:"John Doe"`
	PhoneNumber        string `json:"phone_number" example:"(123) 456-7890"`
	City               string `json:"city" example:"Los Angeles"`
	ZipCode            string `json:"zip_code" example:"90001"`
	DotNumber          string `json:"dot_number" example:"DOT1234567"`
}

type TruckResponse struct {
	ID             int64    `json:"id"`
	TruckName      string   `json:"truck_name"`
	LicensePlate   string   `json:"license_plate"`
	Make           string   `json:"make"`
	Model          string   `json:"model"`
	Year           int      `json:"year"`
	Color          string   `json:"color"`
	Length         float64  `json:"length"`
	Width          float64  `json:"width"`
	Height         float64  `json:"height"`
	MaxWeight      float64  `json:"max_weight"`
	TruckType      string   `json:"truck_type"`
	ClimateControl bool     `json:"climate_control"`
	Liftgate       bool     `json:"liftgate"`
	PalletJack     bool     `json:"pallet_jack"`
	SecuritySystem bool     `json:"security_system"`
	Refrigerated   bool     `json:"refrigerated"`
	FurniturePads  bool     `json:"furniture_pads"`
	PhotoURLs      []string `json:"photo_urls"`
}
