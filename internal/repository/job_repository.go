package repository

import (
	"context"
	"fmt"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (
			contractor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items, 
			inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			truck_size, pickup_address, pickup_floor, pickup_building_type, pickup_walk_distance,
			delivery_address, delivery_floor, delivery_building_type, delivery_walk_distance,
			distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			weight_lbs, volume_cu_ft
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		ctx,
		query,
		job.ContractorID, job.JobType, job.NumberOfBedrooms, job.PackingBoxes, job.BulkyItems,
		job.InventoryList, job.Hoisting, job.AdditionalServicesDescription, job.EstimatedCrewAssistants,
		job.TruckSize, job.PickupAddress, job.PickupFloor, job.PickupBuildingType, job.PickupWalkDistance,
		job.DeliveryAddress, job.DeliveryFloor, job.DeliveryBuildingType, job.DeliveryWalkDistance,
		job.DistanceMiles, "pending", job.PickupDate, job.PickupTimeFrom, job.PickupTimeTo,
		job.DeliveryDate, job.DeliveryTimeFrom, job.DeliveryTimeTo, job.CutAmount, job.PaymentAmount,
		job.WeightLbs, job.VolumeCuFt,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
}

func (r *JobRepository) GetAvailableJobs(ctx context.Context, userID int64, offset, limit int) ([]models.AvailableJobDTO, error) {
	query := `
		SELECT id, job_type, distance_miles, pickup_address, delivery_address,
			   pickup_date, truck_size, weight_lbs, volume_cu_ft, payment_amount,
			   contractor_id
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.AvailableJobDTO
	for rows.Next() {
		var job models.AvailableJobDTO
		err := rows.Scan(
			&job.ID, &job.JobType, &job.DistanceMiles, &job.PickupAddress,
			&job.DeliveryAddress, &job.PickupDate, &job.TruckSize,
			&job.WeightLbs, &job.VolumeCuFt, &job.PaymentAmount,
			&job.ContractorID,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) GetJobByID(ctx context.Context, jobID int64) (*models.Job, error) {
	query := `
		SELECT id, contractor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items,
			   inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			   truck_size, pickup_address, pickup_floor, pickup_building_type, pickup_walk_distance,
			   delivery_address, delivery_floor, delivery_building_type, delivery_walk_distance,
			   distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			   delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			   weight_lbs, volume_cu_ft, created_at, updated_at
		FROM jobs WHERE id = $1`

	var job models.Job
	err := r.db.QueryRow(ctx, query, jobID).Scan(
		&job.ID, &job.ContractorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
		&job.BulkyItems, &job.InventoryList, &job.Hoisting, &job.AdditionalServicesDescription,
		&job.EstimatedCrewAssistants, &job.TruckSize, &job.PickupAddress, &job.PickupFloor,
		&job.PickupBuildingType, &job.PickupWalkDistance, &job.DeliveryAddress, &job.DeliveryFloor,
		&job.DeliveryBuildingType, &job.DeliveryWalkDistance, &job.DistanceMiles, &job.JobStatus,
		&job.PickupDate, &job.PickupTimeFrom, &job.PickupTimeTo, &job.DeliveryDate,
		&job.DeliveryTimeFrom, &job.DeliveryTimeTo, &job.CutAmount, &job.PaymentAmount,
		&job.WeightLbs, &job.VolumeCuFt, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobRepository) DeleteJob(ctx context.Context, jobID, userID int64) error {
	query := `DELETE FROM jobs WHERE id = $1 AND contractor_id = $2`
	result, err := r.db.Exec(ctx, query, jobID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found or you don't have permission to delete it")
	}

	return nil
}

func (r *JobRepository) ClaimJob(ctx context.Context, jobID, userID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var contractorID int64
	var status string
	err = tx.QueryRow(ctx, "SELECT contractor_id, job_status FROM jobs WHERE id = $1", jobID).Scan(&contractorID, &status)
	if err != nil {
		return err
	}

	if contractorID == userID {
		return fmt.Errorf("you cannot claim your own job")
	}

	if status != "active" {
		return fmt.Errorf("job is not available for claiming")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO job_applications (job_id, user_id) 
		VALUES ($1, $2) 
		ON CONFLICT (job_id, user_id) DO NOTHING`,
		jobID, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *JobRepository) GetCountAvailableJobs(ctx context.Context, userID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jobs WHERE contractor_id != $1 AND job_status = 'active'`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *JobRepository) GetMyJobs(ctx context.Context, userID int64, offset, limit int) ([]models.Job, error) {
	query := `
		SELECT id, contractor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items,
			   inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			   truck_size, pickup_address, pickup_floor, pickup_building_type, pickup_walk_distance,
			   delivery_address, delivery_floor, delivery_building_type, delivery_walk_distance,
			   distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			   delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			   weight_lbs, volume_cu_ft, created_at, updated_at
		FROM jobs 
		WHERE contractor_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.ContractorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
			&job.BulkyItems, &job.InventoryList, &job.Hoisting, &job.AdditionalServicesDescription,
			&job.EstimatedCrewAssistants, &job.TruckSize, &job.PickupAddress, &job.PickupFloor,
			&job.PickupBuildingType, &job.PickupWalkDistance, &job.DeliveryAddress, &job.DeliveryFloor,
			&job.DeliveryBuildingType, &job.DeliveryWalkDistance, &job.DistanceMiles, &job.JobStatus,
			&job.PickupDate, &job.PickupTimeFrom, &job.PickupTimeTo, &job.DeliveryDate,
			&job.DeliveryTimeFrom, &job.DeliveryTimeTo, &job.CutAmount, &job.PaymentAmount,
			&job.WeightLbs, &job.VolumeCuFt, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) GetCountMyJobs(ctx context.Context, userID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jobs WHERE contractor_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *JobRepository) JobExists(ctx context.Context, jobID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM jobs 
			WHERE id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, jobID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *JobRepository) GetClaimedJobs(ctx context.Context, userID int64, offset, limit int) ([]models.Job, error) {
	query := `
		SELECT j.id, j.contractor_id, j.job_type, j.number_of_bedrooms, j.packing_boxes, j.bulky_items,
			   j.inventory_list, j.hoisting, j.additional_services_description, j.estimated_crew_assistants,
			   j.truck_size, j.pickup_address, j.pickup_floor, j.pickup_building_type, j.pickup_walk_distance,
			   j.delivery_address, j.delivery_floor, j.delivery_building_type, j.delivery_walk_distance,
			   j.distance_miles, j.job_status, j.pickup_date, j.pickup_time_from, j.pickup_time_to,
			   j.delivery_date, j.delivery_time_from, j.delivery_time_to, j.cut_amount, j.payment_amount,
			   j.weight_lbs, j.volume_cu_ft, j.created_at, j.updated_at
		FROM jobs j
		INNER JOIN job_applications ja ON j.id = ja.job_id
		WHERE ja.user_id = $1
		ORDER BY ja.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.ContractorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
			&job.BulkyItems, &job.InventoryList, &job.Hoisting, &job.AdditionalServicesDescription,
			&job.EstimatedCrewAssistants, &job.TruckSize, &job.PickupAddress, &job.PickupFloor,
			&job.PickupBuildingType, &job.PickupWalkDistance, &job.DeliveryAddress, &job.DeliveryFloor,
			&job.DeliveryBuildingType, &job.DeliveryWalkDistance, &job.DistanceMiles, &job.JobStatus,
			&job.PickupDate, &job.PickupTimeFrom, &job.PickupTimeTo, &job.DeliveryDate,
			&job.DeliveryTimeFrom, &job.DeliveryTimeTo, &job.CutAmount, &job.PaymentAmount,
			&job.WeightLbs, &job.VolumeCuFt, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) GetCountClaimedJobs(ctx context.Context, userID int64) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM jobs j
		INNER JOIN job_applications ja ON j.id = ja.job_id
		WHERE ja.user_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
