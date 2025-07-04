package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Job represents a job entity
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
}

// JobRepository defines the interface for job data operations
type JobRepository interface {
	CreateJob(ctx context.Context, job *Job) error
	GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]Job, error)
	GetUserJobs(ctx context.Context, userID int64) ([]Job, error)
	DeleteJob(ctx context.Context, userID, jobID int64) error
	ApplyForJob(ctx context.Context, userID, jobID int64) error
	GetMyApplications(ctx context.Context, userID int64) ([]Job, error) // New method

}

// jobRepository implements JobRepository
type jobRepository struct {
	db *pgxpool.Pool
}

// NewJobRepository creates a new JobRepository
func NewJobRepository(db *pgxpool.Pool) JobRepository {
	return &jobRepository{db: db}
}

// CreateJob creates a new job in the database
func (r *jobRepository) CreateJob(ctx context.Context, job *Job) error {
	query := `
		INSERT INTO jobs (
			user_id, job_title, description, cargo_type, urgency, truck_size, loading_assistance,
			pickup_date, pickup_time_window, delivery_date, delivery_time_window, pickup_location,
			delivery_location, payout_amount, early_delivery_bonus, payment_terms, weight_lb,
			volume_cu_ft, liftgate, fragile_items, climate_control, assembly_required,
			extra_insurance, additional_packing
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		job.UserID, job.JobTitle, job.Description, job.CargoType, job.Urgency, job.TruckSize,
		job.LoadingAssistance, job.PickupDate, job.PickupTimeWindow, job.DeliveryDate,
		job.DeliveryTimeWindow, job.PickupLocation, job.DeliveryLocation, job.PayoutAmount,
		job.EarlyDeliveryBonus, job.PaymentTerms, job.WeightLb, job.VolumeCuFt, job.Liftgate,
		job.FragileItems, job.ClimateControl, job.AssemblyRequired, job.ExtraInsurance,
		job.AdditionalPacking).Scan(&job.ID)
}

// GetAvailableJobs fetches jobs excluding those created by the given userID with optional filters and pagination
func (r *jobRepository) GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]Job, error) {
	query := `
		SELECT id, user_id, job_title, description, cargo_type, urgency, truck_size, loading_assistance,
		       pickup_date, pickup_time_window, delivery_date, delivery_time_window, pickup_location,
		       delivery_location, payout_amount, early_delivery_bonus, payment_terms, weight_lb,
		       volume_cu_ft, liftgate, fragile_items, climate_control, assembly_required,
		       extra_insurance, additional_packing
		FROM jobs
		WHERE user_id != $1
	`

	params := []interface{}{userID}
	paramCount := 2

	// Apply filters
	if pickupLoc, ok := filters["pickup_location"]; ok && pickupLoc != "" {
		query += fmt.Sprintf(" AND pickup_location ILIKE $%d", paramCount)
		params = append(params, "%"+pickupLoc+"%")
		paramCount++
	}
	if deliveryLoc, ok := filters["delivery_location"]; ok && deliveryLoc != "" {
		query += fmt.Sprintf(" AND delivery_location ILIKE $%d", paramCount)
		params = append(params, "%"+deliveryLoc+"%")
		paramCount++
	}
	if startDate, ok := filters["pickup_date_start"]; ok && startDate != "" {
		query += fmt.Sprintf(" AND pickup_date >= $%d", paramCount)
		t, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, err
		}
		params = append(params, t)
		paramCount++
	}
	if endDate, ok := filters["pickup_date_end"]; ok && endDate != "" {
		query += fmt.Sprintf(" AND pickup_date <= $%d", paramCount)
		t, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, err
		}
		params = append(params, t)
		paramCount++
	}
	if truckSize, ok := filters["truck_size"]; ok && truckSize != "" {
		query += fmt.Sprintf(" AND truck_size = $%d", paramCount)
		params = append(params, truckSize)
		paramCount++
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	params = append(params, limit, offset)

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID, &job.UserID, &job.JobTitle, &job.Description, &job.CargoType, &job.Urgency,
			&job.TruckSize, &job.LoadingAssistance, &job.PickupDate, &job.PickupTimeWindow,
			&job.DeliveryDate, &job.DeliveryTimeWindow, &job.PickupLocation, &job.DeliveryLocation,
			&job.PayoutAmount, &job.EarlyDeliveryBonus, &job.PaymentTerms, &job.WeightLb,
			&job.VolumeCuFt, &job.Liftgate, &job.FragileItems, &job.ClimateControl,
			&job.AssemblyRequired, &job.ExtraInsurance, &job.AdditionalPacking,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetUserJobs fetches all jobs created by the given userID
func (r *jobRepository) GetUserJobs(ctx context.Context, userID int64) ([]Job, error) {
	query := `
		SELECT id, user_id, job_title, description, cargo_type, urgency, truck_size, loading_assistance,
		       pickup_date, pickup_time_window, delivery_date, delivery_time_window, pickup_location,
		       delivery_location, payout_amount, early_delivery_bonus, payment_terms, weight_lb,
		       volume_cu_ft, liftgate, fragile_items, climate_control, assembly_required,
		       extra_insurance, additional_packing
		FROM jobs
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID, &job.UserID, &job.JobTitle, &job.Description, &job.CargoType, &job.Urgency,
			&job.TruckSize, &job.LoadingAssistance, &job.PickupDate, &job.PickupTimeWindow,
			&job.DeliveryDate, &job.DeliveryTimeWindow, &job.PickupLocation, &job.DeliveryLocation,
			&job.PayoutAmount, &job.EarlyDeliveryBonus, &job.PaymentTerms, &job.WeightLb,
			&job.VolumeCuFt, &job.Liftgate, &job.FragileItems, &job.ClimateControl,
			&job.AssemblyRequired, &job.ExtraInsurance, &job.AdditionalPacking,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// DeleteJob deletes a job if it belongs to the given userID
func (r *jobRepository) DeleteJob(ctx context.Context, userID, jobID int64) error {
	query := `
		DELETE FROM jobs
		WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(ctx, query, jobID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found or unauthorized")
	}
	return nil
}

// ApplyForJob creates an application for a job
func (r *jobRepository) ApplyForJob(ctx context.Context, userID, jobID int64) error {
	query := `
		INSERT INTO job_applications (job_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT ON CONSTRAINT unique_application DO NOTHING
		RETURNING id
	`
	var id int64
	err := r.db.QueryRow(ctx, query, jobID, userID).Scan(&id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return fmt.Errorf("application already exists or job not found")
		}
		return err
	}
	return nil
}

// GetMyApplications fetches all jobs the user has applied for
func (r *jobRepository) GetMyApplications(ctx context.Context, userID int64) ([]Job, error) {
	query := `
		SELECT j.id, j.user_id, j.job_title, j.description, j.cargo_type, j.urgency, j.truck_size, j.loading_assistance,
		       j.pickup_date, j.pickup_time_window, j.delivery_date, j.delivery_time_window, j.pickup_location,
		       j.delivery_location, j.payout_amount, j.early_delivery_bonus, j.payment_terms, j.weight_lb,
		       j.volume_cu_ft, j.liftgate, j.fragile_items, j.climate_control, j.assembly_required,
		       j.extra_insurance, j.additional_packing
		FROM jobs j
		JOIN job_applications ja ON j.id = ja.job_id
		WHERE ja.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID, &job.UserID, &job.JobTitle, &job.Description, &job.CargoType, &job.Urgency,
			&job.TruckSize, &job.LoadingAssistance, &job.PickupDate, &job.PickupTimeWindow,
			&job.DeliveryDate, &job.DeliveryTimeWindow, &job.PickupLocation, &job.DeliveryLocation,
			&job.PayoutAmount, &job.EarlyDeliveryBonus, &job.PaymentTerms, &job.WeightLb,
			&job.VolumeCuFt, &job.Liftgate, &job.FragileItems, &job.ClimateControl,
			&job.AssemblyRequired, &job.ExtraInsurance, &job.AdditionalPacking,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
