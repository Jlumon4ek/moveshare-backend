package models

import (
	"time"
)

type Truck struct {
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
	Photos         []string  `json:"photos"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
