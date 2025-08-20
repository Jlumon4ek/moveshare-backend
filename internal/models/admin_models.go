package models

import (
	"time"
)

type UserCompanyInfo struct {
	ID           int64     `json:"id"`
	CompanyName  string    `json:"company_name"`
	Email        string    `json:"email"`
	TrucksNumber int       `json:"trucks_number"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type JobManagementInfo struct {
	ID          int64     `json:"id"`
	Size        string    `json:"size"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Date        time.Time `json:"date"`
	Payout      float64   `json:"payout"`
	Status      string    `json:"status"`
}

type TruckSwagger struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	TruckName      string    `json:"truck_name"`
	LicensePlate   string    `json:"license_plate"`
	Make           string    `json:"make"`
	Model          string    `json:"model"`
	Year           int       `json:"year"`
	Color          string    `json:"color"`
	Length         float64   `json:"length"`
	Width          float64   `json:"width"`
	Height         float64   `json:"height"`
	MaxWeight      float64   `json:"max_weight"`
	TruckType      string    `json:"truck_type"`
	ClimateControl bool      `json:"climate_control"`
	Liftgate       bool      `json:"liftgate"`
	PalletJack     bool      `json:"pallet_jack"`
	SecuritySystem bool      `json:"security_system"`
	Refrigerated   bool      `json:"refrigerated"`
	FurniturePads  bool      `json:"furniture_pads"`
	PhotoURLs      []string  `json:"photo_urls,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserFullInfo struct {
	User         User                    `json:"user"`
	Company      *Company                `json:"company,omitempty"`
	Trucks       []TruckSwagger          `json:"trucks"`
	Jobs         []Job                   `json:"jobs"`
	Reviews      []Review                `json:"reviews"`
	Payments     []Payment               `json:"payments"`
	Verification []VerificationFile      `json:"verification"`
}

type PaginatedUsersResponse struct {
	Users  []UserCompanyInfo `json:"users"`
	Total  int               `json:"total"`
	Page   int               `json:"page"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

type PaginatedJobsResponse struct {
	Jobs   []JobManagementInfo `json:"jobs"`
	Total  int                 `json:"total"`
	Page   int                 `json:"page"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

type TopCompany struct {
	CompanyName string `json:"company_name"`
	JobsCount   int    `json:"jobs_count"`
}

type BusyRoute struct {
	Route           string `json:"route"`
	PickupAddress   string `json:"pickup_address"`
	DeliveryAddress string `json:"delivery_address"`
	JobsCount       int    `json:"jobs_count"`
}

type PlatformAnalytics struct {
	TopCompanies  []TopCompany `json:"top_companies"`
	BusiestRoutes []BusyRoute  `json:"busiest_routes"`
}

type SystemSettings struct {
	ID                int     `json:"id" db:"id"`
	CommissionRate    float64 `json:"commission_rate" db:"commission_rate"`
	NewUserApproval   string  `json:"new_user_approval" db:"new_user_approval"`
	MinimumPayout     int     `json:"minimum_payout" db:"minimum_payout"`
	JobExpirationDays int     `json:"job_expiration_days" db:"job_expiration_days"`
}
