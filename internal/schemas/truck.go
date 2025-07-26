package schemas

import (
	"moveshare/internal/models"
)

type TruckRequest struct {
	TruckName      string   `json:"truck_name" example:"Blue Thunder"`
	LicensePlate   string   `json:"license_plate" example:"ABC123"`
	Make           string   `json:"make" example:"Ford"`
	Model          string   `json:"model" example:"F-150"`
	Year           int      `json:"year" example:"2020"`
	Color          string   `json:"color" example:"Red"`
	Length         float64  `json:"length" example:"26"`
	Width          float64  `json:"width" example:"8.5"`
	Height         float64  `json:"height" example:"9.5"`
	MaxWeight      float64  `json:"max_weight" example:"10000"`
	TruckType      string   `json:"truck_type" example:"Large"`
	ClimateControl bool     `json:"climate_control" example:"true"`
	Liftgate       bool     `json:"liftgate" example:"true"`
	PalletJack     bool     `json:"pallet_jack" example:"false"`
	SecuritySystem bool     `json:"security_system" example:"true"`
	Refrigerated   bool     `json:"refrigerated" example:"false"`
	FurniturePads  bool     `json:"furniture_pads" example:"true"`
	Photos         []string `json:"photos,omitempty"`
}

type TruckResponse struct {
	Truck *models.Truck `json:"truck"`
}

type TrucksResponse struct {
	Trucks []models.Truck `json:"trucks"`
}
