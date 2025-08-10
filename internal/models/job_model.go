package models

import (
	"time"
)

type Job struct {
	ID           int64 `json:"id" db:"id"`
	ContractorID int64 `json:"contractor_id" db:"contractor_id"`

	// Job details
	JobType          string `json:"job_type" db:"job_type"`
	NumberOfBedrooms string `json:"number_of_bedrooms" db:"number_of_bedrooms"`

	// Additional services
	PackingBoxes                  bool   `json:"packing_boxes" db:"packing_boxes"`
	BulkyItems                    bool   `json:"bulky_items" db:"bulky_items"`
	InventoryList                 bool   `json:"inventory_list" db:"inventory_list"`
	Hoisting                      bool   `json:"hoisting" db:"hoisting"`
	AdditionalServicesDescription string `json:"additional_services_description" db:"additional_services_description"`

	// Crew and truck
	EstimatedCrewAssistants string `json:"estimated_crew_assistants" db:"estimated_crew_assistants"`
	TruckSize               string `json:"truck_size" db:"truck_size"`

	// Pickup location
	PickupAddress      string `json:"pickup_address" db:"pickup_address"`
	PickupFloor        *int   `json:"pickup_floor" db:"pickup_floor"`
	PickupBuildingType string `json:"pickup_building_type" db:"pickup_building_type"`
	PickupWalkDistance string `json:"pickup_walk_distance" db:"pickup_walk_distance"`

	// Delivery location
	DeliveryAddress      string `json:"delivery_address" db:"delivery_address"`
	DeliveryFloor        *int   `json:"delivery_floor" db:"delivery_floor"`
	DeliveryBuildingType string `json:"delivery_building_type" db:"delivery_building_type"`
	DeliveryWalkDistance string `json:"delivery_walk_distance" db:"delivery_walk_distance"`

	// Job info
	DistanceMiles float64 `json:"distance_miles" db:"distance_miles"`
	JobStatus     string  `json:"job_status" db:"job_status"`

	// Schedule
	PickupDate       time.Time `json:"pickup_date" db:"pickup_date"`
	PickupTimeFrom   time.Time `json:"pickup_time_from" db:"pickup_time_from"`
	PickupTimeTo     time.Time `json:"pickup_time_to" db:"pickup_time_to"`
	DeliveryDate     time.Time `json:"delivery_date" db:"delivery_date"`
	DeliveryTimeFrom time.Time `json:"delivery_time_from" db:"delivery_time_from"`
	DeliveryTimeTo   time.Time `json:"delivery_time_to" db:"delivery_time_to"`

	// Payment
	CutAmount     float64 `json:"cut_amount" db:"cut_amount"`
	PaymentAmount float64 `json:"payment_amount" db:"payment_amount"`

	// Load details
	WeightLbs  float64 `json:"weight_lbs" db:"weight_lbs"`
	VolumeCuFt float64 `json:"volume_cu_ft" db:"volume_cu_ft"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type JobApplication struct {
	ID        int64     `json:"id" db:"id"`
	JobID     int64     `json:"job_id" db:"job_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// DTO for available jobs list
type AvailableJobDTO struct {
	ID              int64     `json:"id"`
	ContractorID    int64     `json:"contractor_id"`
	JobType         string    `json:"job_type"`
	DistanceMiles   float64   `json:"distance_miles"`
	PickupAddress   string    `json:"pickup_address"`
	DeliveryAddress string    `json:"delivery_address"`
	PickupDate      time.Time `json:"pickup_date"`
	TruckSize       string    `json:"truck_size"`
	WeightLbs       float64   `json:"weight_lbs"`
	VolumeCuFt      float64   `json:"volume_cu_ft"`
	PaymentAmount   float64   `json:"payment_amount"`
}

// Request DTOs
type CreateJobRequest struct {
	JobType          string `json:"job_type" binding:"required"`
	NumberOfBedrooms string `json:"number_of_bedrooms"`

	PackingBoxes                  bool   `json:"packing_boxes"`
	BulkyItems                    bool   `json:"bulky_items"`
	InventoryList                 bool   `json:"inventory_list"`
	Hoisting                      bool   `json:"hoisting"`
	AdditionalServicesDescription string `json:"additional_services_description"`

	EstimatedCrewAssistants string `json:"estimated_crew_assistants"`
	TruckSize               string `json:"truck_size" binding:"required"`

	PickupAddress      string `json:"pickup_address" binding:"required"`
	PickupFloor        *int   `json:"pickup_floor"`
	PickupBuildingType string `json:"pickup_building_type"`
	PickupWalkDistance string `json:"pickup_walk_distance"`

	DeliveryAddress      string `json:"delivery_address" binding:"required"`
	DeliveryFloor        *int   `json:"delivery_floor"`
	DeliveryBuildingType string `json:"delivery_building_type"`
	DeliveryWalkDistance string `json:"delivery_walk_distance"`

	DistanceMiles float64 `json:"distance_miles"`

	PickupDate       string `json:"pickup_date" binding:"required"`
	PickupTimeFrom   string `json:"pickup_time_from" binding:"required"`
	PickupTimeTo     string `json:"pickup_time_to" binding:"required"`
	DeliveryDate     string `json:"delivery_date" binding:"required"`
	DeliveryTimeFrom string `json:"delivery_time_from" binding:"required"`
	DeliveryTimeTo   string `json:"delivery_time_to" binding:"required"`

	CutAmount     float64 `json:"cut_amount"` // <- ДОБАВЬ ЭТО
	PaymentAmount float64 `json:"payment_amount" binding:"required"`
	WeightLbs     float64 `json:"weight_lbs"`
	VolumeCuFt    float64 `json:"volume_cu_ft"`
}

type PaginationQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}
