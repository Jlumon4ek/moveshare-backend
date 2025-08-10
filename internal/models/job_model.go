// internal/models/job_model.go
package models

import (
	"mime/multipart"
	"time"
)

type Job struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`

	// Job Type & Basic Info
	JobType     string `json:"job_type"`    // "Residential Move", "Office Relocation", "Warehouse Transfer", "Other"
	JobTitle    string `json:"job_title"`   // Custom title
	Description string `json:"description"` // Additional description

	// Residential Move Details
	NumberOfBedrooms string `json:"number_of_bedrooms"` // "1 Bedroom", "2 Bedrooms", etc., "Office"

	// Additional Services
	PackingBoxes           bool   `json:"packing_boxes"`            // Packing Boxes (hourly fee)
	BulkyItems             bool   `json:"bulky_items"`              // Bulky Items (piano, safe) (add. fee)
	InventoryList          bool   `json:"inventory_list"`           // Inventory List
	Hoisting               bool   `json:"hoisting"`                 // Hoisting (add. fee)
	AdditionalServicesDesc string `json:"additional_services_desc"` // Description of additional services

	// Truck & Crew Requirements
	TruckSize      string `json:"truck_size"`      // "Small (15')", "Medium (20+')", "Large (26+')"
	CrewAssistants string `json:"crew_assistants"` // "Driver only", "Driver + 1 Assistant", etc.

	// Pickup Location
	PickupLocation     string `json:"pickup_location"`      // Street address
	PickupType         string `json:"pickup_type"`          // "House", "Stairs", "Elevator"
	PickupWalkDistance string `json:"pickup_walk_distance"` // "< 50 ft", "50-100 ft", "100-200 ft", "> 200 ft"

	// Delivery Location
	DeliveryLocation     string `json:"delivery_location"`      // Street address
	DeliveryType         string `json:"delivery_type"`          // "House", "Stairs", "Elevator"
	DeliveryWalkDistance string `json:"delivery_walk_distance"` // "< 50 ft", "50-100 ft", "100-200 ft", "> 200 ft"

	// Schedule
	PickupDate        time.Time `json:"pickup_date"`
	PickupTimeStart   string    `json:"pickup_time_start"` // "6:00 AM"
	PickupTimeEnd     string    `json:"pickup_time_end"`   // "6:30 AM"
	DeliveryDate      time.Time `json:"delivery_date"`
	DeliveryTimeStart string    `json:"delivery_time_start"` // "9:00 AM"
	DeliveryTimeEnd   string    `json:"delivery_time_end"`   // "8:00 AM"

	// Payment Details
	CutAmount     float64 `json:"cut_amount"`     // CUT ($) - комиссия платформы
	PaymentAmount float64 `json:"payment_amount"` // Payment ($) - оплата подрядчику

	// System fields
	Status     string                  `json:"status"`     // "draft", "pending_payment", "active", "applied", "in_progress", "completed", "cancelled"
	PhotoFiles []*multipart.FileHeader `json:"-"`          // Загружаемые файлы (не JSON)
	PhotoURLs  []string                `json:"photo_urls"` // URL фотографий после загрузки
	CreatedAt  time.Time               `json:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at"`

	// Calculated fields
	DistanceMiles float64 `json:"distance_miles"` // Расстояние между точками
	TotalAmount   float64 `json:"total_amount"`   // CUT + Payment
}

// DTO для создания работы
type CreateJobRequest struct {
	// Job Type & Basic Info
	JobType     string `json:"job_type" binding:"required" example:"Residential Move"`
	JobTitle    string `json:"job_title" binding:"required" example:"3 bedroom house move"`
	Description string `json:"description" example:"Moving from house to apartment, careful with fragile items"`

	// Residential Move Details
	NumberOfBedrooms string `json:"number_of_bedrooms" binding:"required" example:"3 Bedrooms"`

	// Additional Services
	PackingBoxes           bool   `json:"packing_boxes" example:"true"`
	BulkyItems             bool   `json:"bulky_items" example:"false"`
	InventoryList          bool   `json:"inventory_list" example:"true"`
	Hoisting               bool   `json:"hoisting" example:"false"`
	AdditionalServicesDesc string `json:"additional_services_desc" example:"Need extra packing for kitchen items"`

	// Truck & Crew Requirements
	TruckSize      string `json:"truck_size" binding:"required" example:"Large (26+')"`
	CrewAssistants string `json:"crew_assistants" binding:"required" example:"Driver + 2 Assistant"`

	// Pickup Location
	PickupLocation     string `json:"pickup_location" binding:"required" example:"123 Main St, Chicago, IL"`
	PickupType         string `json:"pickup_type" binding:"required" example:"Stairs"`
	PickupWalkDistance string `json:"pickup_walk_distance" binding:"required" example:"100-200 ft"`

	// Delivery Location
	DeliveryLocation     string `json:"delivery_location" binding:"required" example:"456 Oak Ave, Milwaukee, WI"`
	DeliveryType         string `json:"delivery_type" binding:"required" example:"Elevator"`
	DeliveryWalkDistance string `json:"delivery_walk_distance" binding:"required" example:"50-100 ft"`

	// Schedule
	PickupDate        string `json:"pickup_date" binding:"required" example:"2025-08-15"` // YYYY-MM-DD
	PickupTimeStart   string `json:"pickup_time_start" binding:"required" example:"8:00 AM"`
	PickupTimeEnd     string `json:"pickup_time_end" binding:"required" example:"9:00 AM"`
	DeliveryDate      string `json:"delivery_date" binding:"required" example:"2025-08-15"` // YYYY-MM-DD
	DeliveryTimeStart string `json:"delivery_time_start" binding:"required" example:"2:00 PM"`
	DeliveryTimeEnd   string `json:"delivery_time_end" binding:"required" example:"3:00 PM"`

	// Payment Details
	CutAmount     float64 `json:"cut_amount" binding:"required,min=0" example:"300"`      // Комиссия платформы
	PaymentAmount float64 `json:"payment_amount" binding:"required,min=0" example:"1500"` // Оплата подрядчику
}

type CreateJobResponse struct {
	JobID   int64  `json:"job_id" example:"123"`
	Message string `json:"message" example:"Job created successfully"`
	Status  string `json:"status" example:"pending_payment"`
	Success bool   `json:"success" example:"true"`
}

type JobResponse struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	JobType     string `json:"job_type"`
	JobTitle    string `json:"job_title"`
	Description string `json:"description"`

	NumberOfBedrooms string `json:"number_of_bedrooms"`

	PackingBoxes           bool   `json:"packing_boxes"`
	BulkyItems             bool   `json:"bulky_items"`
	InventoryList          bool   `json:"inventory_list"`
	Hoisting               bool   `json:"hoisting"`
	AdditionalServicesDesc string `json:"additional_services_desc"`

	TruckSize      string `json:"truck_size"`
	CrewAssistants string `json:"crew_assistants"`

	PickupLocation     string `json:"pickup_location"`
	PickupType         string `json:"pickup_type"`
	PickupWalkDistance string `json:"pickup_walk_distance"`

	DeliveryLocation     string `json:"delivery_location"`
	DeliveryType         string `json:"delivery_type"`
	DeliveryWalkDistance string `json:"delivery_walk_distance"`

	PickupDate        time.Time `json:"pickup_date"`
	PickupTimeStart   string    `json:"pickup_time_start"`
	PickupTimeEnd     string    `json:"pickup_time_end"`
	DeliveryDate      time.Time `json:"delivery_date"`
	DeliveryTimeStart string    `json:"delivery_time_start"`
	DeliveryTimeEnd   string    `json:"delivery_time_end"`

	CutAmount     float64 `json:"cut_amount"`
	PaymentAmount float64 `json:"payment_amount"`
	TotalAmount   float64 `json:"total_amount"`

	Status        string    `json:"status"`
	PhotoURLs     []string  `json:"photo_urls"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DistanceMiles float64   `json:"distance_miles"`
}
