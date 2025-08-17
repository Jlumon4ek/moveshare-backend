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
