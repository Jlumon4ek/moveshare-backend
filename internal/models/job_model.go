package models

import "time"

type Job struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	JobTitle           string    `json:"job_title"`
	Description        string    `json:"description"`
	CargoType          string    `json:"cargo_type"`
	Urgency            string    `json:"urgency"`
	TruckSize          string    `json:"truck_size"`
	LoadingAssistance  bool      `json:"loading_assistance"`
	PickupDate         time.Time `json:"pickup_date"`
	PickupTimeWindow   string    `json:"pickup_time_window"`
	DeliveryDate       time.Time `json:"delivery_date"`
	DeliveryTimeWindow string    `json:"delivery_time_window"`
	PickupLocation     string    `json:"pickup_location"`
	DeliveryLocation   string    `json:"delivery_location"`
	PayoutAmount       float64   `json:"payout_amount"`
	EarlyDeliveryBonus float64   `json:"early_delivery_bonus"`
	PaymentTerms       string    `json:"payment_terms"`
	WeightLb           float64   `json:"weight_lb"`
	VolumeCuFt         float64   `json:"volume_cu_ft"`
	Liftgate           bool      `json:"liftgate"`
	FragileItems       bool      `json:"fragile_items"`
	ClimateControl     bool      `json:"climate_control"`
	AssemblyRequired   bool      `json:"assembly_required"`
	ExtraInsurance     bool      `json:"extra_insurance"`
	AdditionalPacking  bool      `json:"additional_packing"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
