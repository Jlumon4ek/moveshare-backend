package models

import "time"

type Job struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	JobTitle           string    `json:"job_title"`
	Description        string    `json:"description"`
	CargoType          string    `json:"cargo_type"`
	Urgency            string    `json:"urgency"`
	PickupLocation     string    `json:"pickup_location"`
	DeliveryLocation   string    `json:"delivery_location"`
	WeightLb           float64   `json:"weight_lb"`
	VolumeCuFt         float64   `json:"volume_cu_ft"`
	TruckSize          string    `json:"truck_size"`
	LoadingAssistance  bool      `json:"loading_assistance"`
	PickupDate         time.Time `json:"pickup_date"`
	PickupTimeWindow   string    `json:"pickup_time_window"`
	DeliveryDate       time.Time `json:"delivery_date"`
	DeliveryTimeWindow string    `json:"delivery_time_window"`
	PayoutAmount       float64   `json:"payout_amount"`
	EarlyDeliveryBonus float64   `json:"early_delivery_bonus"`
	PaymentTerms       string    `json:"payment_terms"`
	Liftgate           bool      `json:"liftgate"`
	FragileItems       bool      `json:"fragile_items"`
	ClimateControl     bool      `json:"climate_control"`
	AssemblyRequired   bool      `json:"assembly_required"`
	ExtraInsurance     bool      `json:"extra_insurance"`
	AdditionalPacking  bool      `json:"additional_packing"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DistanceMiles      float64   `json:"distance_miles"`
}

type PostNewJobRequest struct {
	JobTitle           string    `json:"job_title" example:"Office Furniture Delivery"`
	Description        string    `json:"description" example:"Deliver desks and chairs to downtown office."`
	CargoType          string    `json:"cargo_type" example:"Furniture"`
	Urgency            string    `json:"urgency" example:"High"`
	PickupLocation     string    `json:"pickup_location" example:"123 Main St, Los Angeles, CA"`
	DeliveryLocation   string    `json:"delivery_location" example:"456 Market St, San Francisco, CA"`
	WeightLb           float64   `json:"weight_lb" example:"1200"`
	VolumeCuFt         float64   `json:"volume_cu_ft" example:"350"`
	TruckSize          string    `json:"truck_size" example:"Large"`
	LoadingAssistance  bool      `json:"loading_assistance" example:"true"`
	PickupDate         time.Time `json:"pickup_date" example:"2025-07-28T09:00:00Z"`
	PickupTimeWindow   string    `json:"pickup_time_window" example:"09:00-12:00"`
	DeliveryDate       time.Time `json:"delivery_date" example:"2025-07-29T14:00:00Z"`
	DeliveryTimeWindow string    `json:"delivery_time_window" example:"14:00-17:00"`
	PayoutAmount       float64   `json:"payout_amount" example:"450.00"`
	EarlyDeliveryBonus float64   `json:"early_delivery_bonus" example:"50.00"`
	PaymentTerms       string    `json:"payment_terms" example:"Net 7"`
	Liftgate           bool      `json:"liftgate" example:"true"`
	FragileItems       bool      `json:"fragile_items" example:"false"`
	ClimateControl     bool      `json:"climate_control" example:"false"`
	AssemblyRequired   bool      `json:"assembly_required" example:"true"`
	ExtraInsurance     bool      `json:"extra_insurance" example:"true"`
	AdditionalPacking  bool      `json:"additional_packing" example:"false"`
}

type JobResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

type AvailableJobsResponse struct {
	Jobs []Job `json:"jobs"`
}

type MyJobsResponse struct {
	Jobs []Job `json:"jobs"`
}
