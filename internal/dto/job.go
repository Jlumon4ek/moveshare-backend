package dto

import "moveshare/internal/models"

type JobResponse struct {
	ID                 int64   `json:"id"`
	JobTitle           string  `json:"job_title"`
	Description        string  `json:"description"`
	CargoType          string  `json:"cargo_type"`
	Urgency            string  `json:"urgency"`
	PickupLocation     string  `json:"pickup_location"`
	DeliveryLocation   string  `json:"delivery_location"`
	WeightLb           float64 `json:"weight_lb"`
	VolumeCuFt         float64 `json:"volume_cu_ft"`
	TruckSize          string  `json:"truck_size"`
	LoadingAssistance  bool    `json:"loading_assistance"`
	PickupDate         string  `json:"pickup_date"`
	PickupTimeWindow   string  `json:"pickup_time_window"`
	DeliveryDate       string  `json:"delivery_date"`
	DeliveryTimeWindow string  `json:"delivery_time_window"`
	PayoutAmount       float64 `json:"payout_amount"`
	EarlyDeliveryBonus float64 `json:"early_delivery_bonus"`
	PaymentTerms       string  `json:"payment_terms"`
	Liftgate           bool    `json:"liftgate"`
	FragileItems       bool    `json:"fragile_items"`
	ClimateControl     bool    `json:"climate_control"`
	AssemblyRequired   bool    `json:"assembly_required"`
	ExtraInsurance     bool    `json:"extra_insurance"`
	AdditionalPacking  bool    `json:"additional_packing"`
	Status             string  `json:"status"`
	DistanceMiles      float64 `json:"distance_miles"`
}

func NewJobResponse(j models.Job) JobResponse {
	return JobResponse{
		ID:                 j.ID,
		JobTitle:           j.JobTitle,
		Description:        j.Description,
		CargoType:          j.CargoType,
		Urgency:            j.Urgency,
		PickupLocation:     j.PickupLocation,
		DeliveryLocation:   j.DeliveryLocation,
		WeightLb:           j.WeightLb,
		VolumeCuFt:         j.VolumeCuFt,
		TruckSize:          j.TruckSize,
		LoadingAssistance:  j.LoadingAssistance,
		PickupDate:         j.PickupDate.Format("2006-01-02"),
		PickupTimeWindow:   j.PickupTimeWindow,
		DeliveryDate:       j.DeliveryDate.Format("2006-01-02"),
		DeliveryTimeWindow: j.DeliveryTimeWindow,
		PayoutAmount:       j.PayoutAmount,
		EarlyDeliveryBonus: j.EarlyDeliveryBonus,
		PaymentTerms:       j.PaymentTerms,
		Liftgate:           j.Liftgate,
		FragileItems:       j.FragileItems,
		ClimateControl:     j.ClimateControl,
		AssemblyRequired:   j.AssemblyRequired,
		ExtraInsurance:     j.ExtraInsurance,
		AdditionalPacking:  j.AdditionalPacking,
		Status:             j.Status,
		DistanceMiles:      j.DistanceMiles,
	}
}
