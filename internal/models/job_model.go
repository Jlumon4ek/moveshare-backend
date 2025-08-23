package models

import (
	"fmt"
	"strings"
	"time"
)

type Job struct {
	ID           int64  `json:"id" db:"id"`
	ContractorID int64  `json:"contractor_id" db:"contractor_id"`
	ExecutorID   *int64 `json:"executor_id" db:"executor_id"`

	// Job details
	JobType          string `json:"job_type" db:"job_type"`
	NumberOfBedrooms string `json:"number_of_bedrooms" db:"number_of_bedrooms"`

	// Additional services
	PackingBoxes                  bool    `json:"packing_boxes" db:"packing_boxes"`
	BulkyItems                    bool    `json:"bulky_items" db:"bulky_items"`
	InventoryList                 bool    `json:"inventory_list" db:"inventory_list"`
	Hoisting                      bool    `json:"hoisting" db:"hoisting"`
	AdditionalServicesDescription *string `json:"additional_services_description" db:"additional_services_description"`

	// Crew and truck
	EstimatedCrewAssistants string `json:"estimated_crew_assistants" db:"estimated_crew_assistants"`
	TruckSize               string `json:"truck_size" db:"truck_size"`

	// Pickup location
	PickupAddress      string `json:"pickup_address" db:"pickup_address"`
	PickupCity         string `json:"pickup_city" db:"pickup_city"`
	PickupState        string `json:"pickup_state" db:"pickup_state"`
	PickupFloor        *int   `json:"pickup_floor" db:"pickup_floor"`
	PickupBuildingType string `json:"pickup_building_type" db:"pickup_building_type"`
	PickupWalkDistance string `json:"pickup_walk_distance" db:"pickup_walk_distance"`

	// Delivery location
	DeliveryAddress      string `json:"delivery_address" db:"delivery_address"`
	DeliveryCity         string `json:"delivery_city" db:"delivery_city"`
	DeliveryState        string `json:"delivery_state" db:"delivery_state"`
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

	// Files
	Files []JobFile `json:"files,omitempty"`

	// Contractor info (only for detailed job view)
	ContractorUsername *string  `json:"contractor_username,omitempty"`
	ContractorStatus   *string  `json:"contractor_status,omitempty"`
	ContractorRating   *float64 `json:"contractor_rating,omitempty"`

	// Executor info (only for detailed job view)
	ExecutorUsername *string  `json:"executor_username,omitempty"`
	ExecutorName     *string  `json:"executor_name,omitempty"`
	ExecutorRating   *float64 `json:"executor_rating,omitempty"`

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

type JobFile struct {
	ID          int64     `json:"id" db:"id"`
	JobID       int64     `json:"job_id" db:"job_id"`
	FileID      string    `json:"file_id" db:"file_id"`
	FileName    string    `json:"file_name" db:"file_name"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	ContentType string    `json:"content_type" db:"content_type"`
	FileType    string    `json:"file_type" db:"file_type"` // verification_document or work_photo
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
	FileURL     string    `json:"file_url,omitempty"`
}

// Request DTOs
type CreateJobRequest struct {
	JobType          string `json:"job_type" binding:"required"`
	NumberOfBedrooms string `json:"number_of_bedrooms"`

	PackingBoxes                  bool    `json:"packing_boxes"`
	BulkyItems                    bool    `json:"bulky_items"`
	InventoryList                 bool    `json:"inventory_list"`
	Hoisting                      bool    `json:"hoisting"`
	AdditionalServicesDescription *string `json:"additional_services_description"`

	EstimatedCrewAssistants string `json:"estimated_crew_assistants"`
	TruckSize               string `json:"truck_size" binding:"required"`

	PickupAddress      string `json:"pickup_address" binding:"required"`
	PickupCity         string `json:"pickup_city" binding:"required"`
	PickupState        string `json:"pickup_state" binding:"required"`
	PickupFloor        *int   `json:"pickup_floor"`
	PickupBuildingType string `json:"pickup_building_type"`
	PickupWalkDistance string `json:"pickup_walk_distance"`

	DeliveryAddress      string `json:"delivery_address" binding:"required"`
	DeliveryCity         string `json:"delivery_city" binding:"required"`
	DeliveryState        string `json:"delivery_state" binding:"required"`
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

	CutAmount     float64 `json:"cut_amount"`
	PaymentAmount float64 `json:"payment_amount" binding:"required"`
	WeightLbs     float64 `json:"weight_lbs"`
	VolumeCuFt    float64 `json:"volume_cu_ft"`
}

// CreateJobWithPaymentRequest combines job creation with payment processing
type CreateJobWithPaymentRequest struct {
	// Job details
	JobType          string `json:"job_type" binding:"required"`
	NumberOfBedrooms string `json:"number_of_bedrooms"`

	PackingBoxes                  bool    `json:"packing_boxes"`
	BulkyItems                    bool    `json:"bulky_items"`
	InventoryList                 bool    `json:"inventory_list"`
	Hoisting                      bool    `json:"hoisting"`
	AdditionalServicesDescription *string `json:"additional_services_description"`

	EstimatedCrewAssistants string `json:"estimated_crew_assistants"`
	TruckSize               string `json:"truck_size" binding:"required"`

	PickupAddress      string `json:"pickup_address" binding:"required"`
	PickupCity         string `json:"pickup_city" binding:"required"`
	PickupState        string `json:"pickup_state" binding:"required"`
	PickupFloor        *int   `json:"pickup_floor"`
	PickupBuildingType string `json:"pickup_building_type"`
	PickupWalkDistance string `json:"pickup_walk_distance"`

	DeliveryAddress      string `json:"delivery_address" binding:"required"`
	DeliveryCity         string `json:"delivery_city" binding:"required"`
	DeliveryState        string `json:"delivery_state" binding:"required"`
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

	CutAmount     float64 `json:"cut_amount"`
	PaymentAmount float64 `json:"payment_amount" binding:"required"`
	WeightLbs     float64 `json:"weight_lbs"`
	VolumeCuFt    float64 `json:"volume_cu_ft"`

	// Payment information
	PaymentMethodID *int64 `json:"payment_method_id,omitempty"` // Optional, will use default if not provided
}

type PaginationQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

// internal/models/job_model.go - добавить к существующим структурам

// JobFilters представляет параметры фильтрации для поиска заданий
type JobFilters struct {
	// Пагинация
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`

	// Фильтры
	NumberOfBedrooms *string  `form:"number_of_bedrooms"` // например: "1", "2", "3", "4+", "Studio"
	Origin           *string  `form:"pickup_location"`    // pickup city, state
	Destination      *string  `form:"delivery_location"`  // delivery city, state
	MaxDistance      *float64 `form:"max_distance"`       // максимальная дистанция в милях
	DateStart        *string  `form:"pickup_date_start"`  // начальная дата в формате YYYY-MM-DD
	DateEnd          *string  `form:"pickup_date_end"`    // конечная дата в формате YYYY-MM-DD
	TruckSize        *string  `form:"truck_size"`         // размер грузовика: "Small", "Medium", "Large"
	PayoutMin        *float64 `form:"payout_min"`         // минимальная оплата
	PayoutMax        *float64 `form:"payout_max"`         // максимальная оплата
}

// Validate валидирует параметры фильтрации
func (f *JobFilters) Validate() error {
	// Валидация дат
	if f.DateStart != nil {
		if _, err := time.Parse("2006-01-02", *f.DateStart); err != nil {
			return fmt.Errorf("invalid pickup_date_start format, use YYYY-MM-DD")
		}
	}

	if f.DateEnd != nil {
		if _, err := time.Parse("2006-01-02", *f.DateEnd); err != nil {
			return fmt.Errorf("invalid pickup_date_end format, use YYYY-MM-DD")
		}
	}

	// Валидация дистанции
	if f.MaxDistance != nil && *f.MaxDistance <= 0 {
		return fmt.Errorf("max_distance must be greater than 0")
	}

	// Валидация размера грузовика
	if f.TruckSize != nil {
		validSizes := map[string]bool{"Small": true, "Medium": true, "Large": true}
		sizes := strings.Fields(*f.TruckSize)
		for _, size := range sizes {
			if !validSizes[size] {
				return fmt.Errorf("truck_size must contain only: Small, Medium, Large")
			}
		}
	}

	// Валидация диапазона оплаты
	if f.PayoutMin != nil && *f.PayoutMin < 0 {
		return fmt.Errorf("payout_min must be non-negative")
	}

	if f.PayoutMax != nil && *f.PayoutMax < 0 {
		return fmt.Errorf("payout_max must be non-negative")
	}

	if f.PayoutMin != nil && f.PayoutMax != nil && *f.PayoutMin > *f.PayoutMax {
		return fmt.Errorf("payout_min cannot be greater than payout_max")
	}

	return nil
}

// LocationOption представляет опцию локации
type LocationOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// JobFilterOptions представляет доступные опции для фильтрации
type JobFilterOptions struct {
	NumberOfBedrooms []string          `json:"number_of_bedrooms"`
	TruckSizes       []string          `json:"truck_sizes"`
	PickupLocations  []LocationOption  `json:"pickup_locations"`
	DeliveryLocations []LocationOption `json:"delivery_locations"`
	PayoutRange      struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"payout_range"`
	MaxDistance float64 `json:"max_distance"`
	DateRange struct {
		Min string `json:"min"` // YYYY-MM-DD
		Max string `json:"max"` // YYYY-MM-DD
	} `json:"date_range"`
}

// Также обновим AvailableJobDTO чтобы включить number_of_bedrooms
type AvailableJobDTO struct {
	ID               int64     `json:"id"`
	ContractorID     int64     `json:"contractor_id"`
	JobType          string    `json:"job_type"`
	NumberOfBedrooms string    `json:"number_of_bedrooms"`
	DistanceMiles    float64   `json:"distance_miles"`
	PickupAddress    string    `json:"pickup_address"`
	PickupCity       string    `json:"pickup_city"`
	PickupState      string    `json:"pickup_state"`
	DeliveryAddress  string    `json:"delivery_address"`
	DeliveryCity     string    `json:"delivery_city"`
	DeliveryState    string    `json:"delivery_state"`
	PickupDate       time.Time `json:"pickup_date"`
	DeliveryDate     time.Time `json:"delivery_date"`
	TruckSize        string    `json:"truck_size"`
	WeightLbs        float64   `json:"weight_lbs"`
	VolumeCuFt       float64   `json:"volume_cu_ft"`
	PaymentAmount    float64   `json:"payment_amount"`
	CutAmount        float64   `json:"cut_amount"`
}

// ExportJobsRequest представляет запрос на экспорт работ
type ExportJobsRequest struct {
	JobIDs []int64 `json:"job_ids" binding:"required,min=1"`
}

// CancelJobsRequest представляет запрос на отмену работ
type CancelJobsRequest struct {
	JobIDs []int64 `json:"job_ids" binding:"required,min=1"`
}

// JobsStats представляет статистику по работам
type JobsStats struct {
	ActiveJobsCount    int            `json:"active_jobs_count"`
	NewJobsThisWeek    int            `json:"new_jobs_this_week"`
	StatusDistribution map[string]int `json:"status_distribution"`
}

// UserWorkStats представляет статистику работ пользователя (на которые он откликался)
type UserWorkStats struct {
	CompletedJobs int     `json:"completed_jobs"`
	Earnings      float64 `json:"earnings"`
	UpcomingJobs  int     `json:"upcoming_jobs"`
}
