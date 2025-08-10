package repository

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"strings"
	"time"

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

	_, err = tx.Exec(ctx, `
		UPDATE jobs 
		SET job_status = 'pending', updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1`,
		jobID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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

// internal/repository/job_repository.go - обновить метод GetAvailableJobs

func (r *JobRepository) GetAvailableJobs(ctx context.Context, userID int64, filters *models.JobFilters) ([]models.AvailableJobDTO, int, error) {
	offset := (filters.Page - 1) * filters.Limit

	// Базовый запрос
	baseQuery := `
		SELECT id, job_type, distance_miles, pickup_address, delivery_address,
			   pickup_date, truck_size, weight_lbs, volume_cu_ft, payment_amount,
			   contractor_id, number_of_bedrooms
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
	`

	countQuery := `
		SELECT COUNT(*) 
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
	`

	// Массивы для условий и параметров
	var conditions []string
	var params []interface{}
	params = append(params, userID) // $1
	paramIndex := 2

	// Добавляем фильтры
	if filters.NumberOfBedrooms != nil && *filters.NumberOfBedrooms != "" {
		conditions = append(conditions, fmt.Sprintf("number_of_bedrooms = $%d", paramIndex))
		params = append(params, *filters.NumberOfBedrooms)
		paramIndex++
	}

	if filters.Origin != nil && *filters.Origin != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_address ILIKE $%d", paramIndex))
		params = append(params, "%"+*filters.Origin+"%")
		paramIndex++
	}

	if filters.Destination != nil && *filters.Destination != "" {
		conditions = append(conditions, fmt.Sprintf("delivery_address ILIKE $%d", paramIndex))
		params = append(params, "%"+*filters.Destination+"%")
		paramIndex++
	}

	if filters.MaxDistance != nil {
		conditions = append(conditions, fmt.Sprintf("distance_miles <= $%d", paramIndex))
		params = append(params, *filters.MaxDistance)
		paramIndex++
	}

	if filters.DateStart != nil && *filters.DateStart != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_date >= $%d", paramIndex))
		params = append(params, *filters.DateStart)
		paramIndex++
	}

	if filters.DateEnd != nil && *filters.DateEnd != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_date <= $%d", paramIndex))
		params = append(params, *filters.DateEnd)
		paramIndex++
	}

	if filters.TruckSize != nil && *filters.TruckSize != "" {
		conditions = append(conditions, fmt.Sprintf("truck_size = $%d", paramIndex))
		params = append(params, *filters.TruckSize)
		paramIndex++
	}

	if filters.PayoutMin != nil {
		conditions = append(conditions, fmt.Sprintf("payment_amount >= $%d", paramIndex))
		params = append(params, *filters.PayoutMin)
		paramIndex++
	}

	if filters.PayoutMax != nil {
		conditions = append(conditions, fmt.Sprintf("payment_amount <= $%d", paramIndex))
		params = append(params, *filters.PayoutMax)
		paramIndex++
	}

	// Добавляем условия к запросам
	if len(conditions) > 0 {
		conditionStr := " AND " + strings.Join(conditions, " AND ")
		baseQuery += conditionStr
		countQuery += conditionStr
	}

	// Получаем общее количество
	var total int
	err := r.db.QueryRow(ctx, countQuery, params...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	// Добавляем сортировку и пагинацию
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	params = append(params, filters.Limit, offset)

	// Выполняем запрос
	rows, err := r.db.Query(ctx, baseQuery, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.AvailableJobDTO
	for rows.Next() {
		var job models.AvailableJobDTO
		err := rows.Scan(
			&job.ID, &job.JobType, &job.DistanceMiles, &job.PickupAddress,
			&job.DeliveryAddress, &job.PickupDate, &job.TruckSize,
			&job.WeightLbs, &job.VolumeCuFt, &job.PaymentAmount,
			&job.ContractorID, &job.NumberOfBedrooms,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return jobs, total, nil
}

// Также нужно обновить GetCountAvailableJobs для работы с фильтрами
func (r *JobRepository) GetCountAvailableJobsWithFilters(ctx context.Context, userID int64, filters *models.JobFilters) (int, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
	`

	var conditions []string
	var params []interface{}
	params = append(params, userID)
	paramIndex := 2

	// Добавляем те же фильтры, что и в основном запросе
	if filters.NumberOfBedrooms != nil && *filters.NumberOfBedrooms != "" {
		conditions = append(conditions, fmt.Sprintf("number_of_bedrooms = $%d", paramIndex))
		params = append(params, *filters.NumberOfBedrooms)
		paramIndex++
	}

	if filters.Origin != nil && *filters.Origin != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_address ILIKE $%d", paramIndex))
		params = append(params, "%"+*filters.Origin+"%")
		paramIndex++
	}

	if filters.Destination != nil && *filters.Destination != "" {
		conditions = append(conditions, fmt.Sprintf("delivery_address ILIKE $%d", paramIndex))
		params = append(params, "%"+*filters.Destination+"%")
		paramIndex++
	}

	if filters.MaxDistance != nil {
		conditions = append(conditions, fmt.Sprintf("distance_miles <= $%d", paramIndex))
		params = append(params, *filters.MaxDistance)
		paramIndex++
	}

	if filters.DateStart != nil && *filters.DateStart != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_date >= $%d", paramIndex))
		params = append(params, *filters.DateStart)
		paramIndex++
	}

	if filters.DateEnd != nil && *filters.DateEnd != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_date <= $%d", paramIndex))
		params = append(params, *filters.DateEnd)
		paramIndex++
	}

	if filters.TruckSize != nil && *filters.TruckSize != "" {
		conditions = append(conditions, fmt.Sprintf("truck_size = $%d", paramIndex))
		params = append(params, *filters.TruckSize)
		paramIndex++
	}

	if filters.PayoutMin != nil {
		conditions = append(conditions, fmt.Sprintf("payment_amount >= $%d", paramIndex))
		params = append(params, *filters.PayoutMin)
		paramIndex++
	}

	if filters.PayoutMax != nil {
		conditions = append(conditions, fmt.Sprintf("payment_amount <= $%d", paramIndex))
		params = append(params, *filters.PayoutMax)
		paramIndex++
	}

	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRow(ctx, countQuery, params...).Scan(&count)
	return count, err
}

func (r *JobRepository) GetCountAvailableJobs(ctx context.Context, userID int64, filters *models.JobFilters) (int, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
	`

	var conditions []string
	var params []interface{}
	params = append(params, userID)
	paramIndex := 2

	if filters.NumberOfBedrooms != nil && *filters.NumberOfBedrooms != "" {
		conditions = append(conditions, fmt.Sprintf("number_of_bedrooms = $%d", paramIndex))
		params = append(params, *filters.NumberOfBedrooms)
		paramIndex++
	}

	if filters.Origin != nil && *filters.Origin != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_address ILIKE $%d", paramIndex))
		params = append(params, "%"+*filters.Origin+"%")
		paramIndex++
	}

	if filters.Destination != nil && *filters.Destination != "" {
		conditions = append(conditions, fmt.Sprintf("delivery_address ILIKE $%d", paramIndex))
		params = append(params, "%"+*filters.Destination+"%")
		paramIndex++
	}

	if filters.MaxDistance != nil {
		conditions = append(conditions, fmt.Sprintf("distance_miles <= $%d", paramIndex))
		params = append(params, *filters.MaxDistance)
		paramIndex++
	}

	if filters.DateStart != nil && *filters.DateStart != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_date >= $%d", paramIndex))
		params = append(params, *filters.DateStart)
		paramIndex++
	}

	if filters.DateEnd != nil && *filters.DateEnd != "" {
		conditions = append(conditions, fmt.Sprintf("pickup_date <= $%d", paramIndex))
		params = append(params, *filters.DateEnd)
		paramIndex++
	}

	if filters.TruckSize != nil && *filters.TruckSize != "" {
		conditions = append(conditions, fmt.Sprintf("truck_size = $%d", paramIndex))
		params = append(params, *filters.TruckSize)
		paramIndex++
	}

	if filters.PayoutMin != nil {
		conditions = append(conditions, fmt.Sprintf("payment_amount >= $%d", paramIndex))
		params = append(params, *filters.PayoutMin)
		paramIndex++
	}

	if filters.PayoutMax != nil {
		conditions = append(conditions, fmt.Sprintf("payment_amount <= $%d", paramIndex))
		params = append(params, *filters.PayoutMax)
		paramIndex++
	}

	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRow(ctx, countQuery, params...).Scan(&count)
	return count, err
}

// internal/repository/job_repository.go - добавить метод GetFilterOptions

func (r *JobRepository) GetFilterOptions(ctx context.Context, userID int64) (*models.JobFilterOptions, error) {
	options := &models.JobFilterOptions{}

	// Получаем уникальные значения количества спален
	bedroomsQuery := `
		SELECT DISTINCT number_of_bedrooms 
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active' 
		AND number_of_bedrooms IS NOT NULL AND number_of_bedrooms != ''
		ORDER BY number_of_bedrooms
	`

	rows, err := r.db.Query(ctx, bedroomsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bedrooms options: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bedroom string
		if err := rows.Scan(&bedroom); err != nil {
			return nil, err
		}
		options.NumberOfBedrooms = append(options.NumberOfBedrooms, bedroom)
	}

	// Получаем уникальные размеры грузовиков
	truckSizesQuery := `
		SELECT DISTINCT truck_size 
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active' 
		AND truck_size IS NOT NULL AND truck_size != ''
		ORDER BY truck_size
	`

	rows, err = r.db.Query(ctx, truckSizesQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get truck sizes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var truckSize string
		if err := rows.Scan(&truckSize); err != nil {
			return nil, err
		}
		options.TruckSizes = append(options.TruckSizes, truckSize)
	}

	// Получаем диапазон оплаты
	payoutRangeQuery := `
		SELECT MIN(payment_amount), MAX(payment_amount)
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
	`

	err = r.db.QueryRow(ctx, payoutRangeQuery, userID).Scan(
		&options.PayoutRange.Min,
		&options.PayoutRange.Max,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payout range: %w", err)
	}

	// Получаем максимальную дистанцию (для слайдера)
	maxDistanceQuery := `
		SELECT MAX(distance_miles)
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
	`

	err = r.db.QueryRow(ctx, maxDistanceQuery, userID).Scan(&options.MaxDistance.Max)
	if err != nil {
		return nil, fmt.Errorf("failed to get max distance: %w", err)
	}

	// Получаем диапазон дат
	dateRangeQuery := `
		SELECT MIN(pickup_date), MAX(pickup_date)
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active'
	`

	var minDate, maxDate time.Time
	err = r.db.QueryRow(ctx, dateRangeQuery, userID).Scan(&minDate, &maxDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get date range: %w", err)
	}

	options.DateRange.Min = minDate.Format("2006-01-02")
	options.DateRange.Max = maxDate.Format("2006-01-02")

	return options, nil
}

func (r *JobRepository) MarkJobCompleted(ctx context.Context, jobID, userID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 
			FROM job_applications 
			WHERE job_id = $1 AND user_id = $2
		)`, jobID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("you have not applied for this job or it doesn't exist")
	}

	var currentStatus string
	err = tx.QueryRow(ctx, "SELECT job_status FROM jobs WHERE id = $1", jobID).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus == "completed" {
		return fmt.Errorf("job is already completed")
	}

	if currentStatus != "pending" && currentStatus != "active" {
		return fmt.Errorf("job status must be pending or active to mark as completed")
	}

	_, err = tx.Exec(ctx, `
		UPDATE jobs 
		SET job_status = 'completed', updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1`,
		jobID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *JobRepository) GetJobsByIDs(ctx context.Context, userID int64, jobIDs []int64) ([]models.Job, error) {
	if len(jobIDs) == 0 {
		return []models.Job{}, nil
	}

	query := `
		SELECT id, contractor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items,
			   inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			   truck_size, pickup_address, pickup_floor, pickup_building_type, pickup_walk_distance,
			   delivery_address, delivery_floor, delivery_building_type, delivery_walk_distance,
			   distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			   delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			   weight_lbs, volume_cu_ft, created_at, updated_at
		FROM jobs 
		WHERE contractor_id = $1 AND id = ANY($2)
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID, jobIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs by IDs: %w", err)
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
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return jobs, nil
}

func (r *JobRepository) GetJobsStats(ctx context.Context, userID int64) (models.JobsStats, error) {
	var stats models.JobsStats
	stats.StatusDistribution = make(map[string]int)

	// Get active jobs count created by this user
	activeJobsQuery := `SELECT COUNT(*) FROM jobs WHERE job_status = 'active' AND contractor_id = $1`
	err := r.db.QueryRow(ctx, activeJobsQuery, userID).Scan(&stats.ActiveJobsCount)
	if err != nil {
		return stats, err
	}

	// Get new jobs this week created by this user (created in the last 7 days)
	newJobsQuery := `SELECT COUNT(*) FROM jobs WHERE created_at >= NOW() - INTERVAL '7 days' AND contractor_id = $1`
	err = r.db.QueryRow(ctx, newJobsQuery, userID).Scan(&stats.NewJobsThisWeek)
	if err != nil {
		return stats, err
	}

	// Get status distribution for user's jobs
	statusQuery := `
		SELECT job_status, COUNT(*) 
		FROM jobs 
		WHERE contractor_id = $1 
		GROUP BY job_status
	`
	rows, err := r.db.Query(ctx, statusQuery, userID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return stats, err
		}
		stats.StatusDistribution[status] = count
	}

	if err = rows.Err(); err != nil {
		return stats, err
	}

	return stats, nil
}
