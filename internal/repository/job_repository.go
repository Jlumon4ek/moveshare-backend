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

// Helper function to scan job data with new city/state fields
func scanJob(scanner interface {
	Scan(dest ...interface{}) error
}, job *models.Job) error {
	return scanner.Scan(
		&job.ID, &job.ContractorID, &job.ExecutorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
		&job.BulkyItems, &job.InventoryList, &job.Hoisting, &job.AdditionalServicesDescription,
		&job.EstimatedCrewAssistants, &job.TruckSize, &job.PickupAddress, &job.PickupCity, &job.PickupState, &job.PickupFloor,
		&job.PickupBuildingType, &job.PickupWalkDistance, &job.DeliveryAddress, &job.DeliveryCity, &job.DeliveryState, &job.DeliveryFloor,
		&job.DeliveryBuildingType, &job.DeliveryWalkDistance, &job.DistanceMiles, &job.JobStatus,
		&job.PickupDate, &job.PickupTimeFrom, &job.PickupTimeTo, &job.DeliveryDate,
		&job.DeliveryTimeFrom, &job.DeliveryTimeTo, &job.CutAmount, &job.PaymentAmount,
		&job.WeightLbs, &job.VolumeCuFt, &job.CreatedAt, &job.UpdatedAt,
	)
}

func (r *JobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (
			contractor_id, executor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items, 
			inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			truck_size, pickup_address, pickup_city, pickup_state, pickup_floor, pickup_building_type, pickup_walk_distance,
			delivery_address, delivery_city, delivery_state, delivery_floor, delivery_building_type, delivery_walk_distance,
			distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			weight_lbs, volume_cu_ft
		) VALUES (
			$1, NULL, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34
		) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		ctx,
		query,
		job.ContractorID, job.JobType, job.NumberOfBedrooms, job.PackingBoxes, job.BulkyItems,
		job.InventoryList, job.Hoisting, job.AdditionalServicesDescription, job.EstimatedCrewAssistants,
		job.TruckSize, job.PickupAddress, job.PickupCity, job.PickupState, job.PickupFloor, job.PickupBuildingType, job.PickupWalkDistance,
		job.DeliveryAddress, job.DeliveryCity, job.DeliveryState, job.DeliveryFloor, job.DeliveryBuildingType, job.DeliveryWalkDistance,
		job.DistanceMiles, job.JobStatus, job.PickupDate, job.PickupTimeFrom, job.PickupTimeTo,
		job.DeliveryDate, job.DeliveryTimeFrom, job.DeliveryTimeTo, job.CutAmount, job.PaymentAmount,
		job.WeightLbs, job.VolumeCuFt,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
}

func (r *JobRepository) GetJobByID(ctx context.Context, jobID int64) (*models.Job, error) {
	query := `
		SELECT j.id, j.contractor_id, j.executor_id, j.job_type, j.number_of_bedrooms, j.packing_boxes, j.bulky_items,
			   j.inventory_list, j.hoisting, j.additional_services_description, j.estimated_crew_assistants,
			   j.truck_size, j.pickup_address, j.pickup_city, j.pickup_state, j.pickup_floor, j.pickup_building_type, j.pickup_walk_distance,
			   j.delivery_address, j.delivery_city, j.delivery_state, j.delivery_floor, j.delivery_building_type, j.delivery_walk_distance,
			   j.distance_miles, j.job_status, j.pickup_date, j.pickup_time_from, j.pickup_time_to,
			   j.delivery_date, j.delivery_time_from, j.delivery_time_to, j.cut_amount, j.payment_amount,
			   j.weight_lbs, j.volume_cu_ft, j.created_at, j.updated_at,
			   u.username, u.status, 
			   COALESCE(AVG(r.rating), 0) as avg_rating
		FROM jobs j
		LEFT JOIN users u ON j.contractor_id = u.id
		LEFT JOIN reviews r ON r.reviewee_id = u.id
		WHERE j.id = $1
		GROUP BY j.id, j.contractor_id, j.executor_id, j.job_type, j.number_of_bedrooms, j.packing_boxes, j.bulky_items,
				 j.inventory_list, j.hoisting, j.additional_services_description, j.estimated_crew_assistants,
				 j.truck_size, j.pickup_address, j.pickup_city, j.pickup_state, j.pickup_floor, j.pickup_building_type, j.pickup_walk_distance,
				 j.delivery_address, j.delivery_city, j.delivery_state, j.delivery_floor, j.delivery_building_type, j.delivery_walk_distance,
				 j.distance_miles, j.job_status, j.pickup_date, j.pickup_time_from, j.pickup_time_to,
				 j.delivery_date, j.delivery_time_from, j.delivery_time_to, j.cut_amount, j.payment_amount,
				 j.weight_lbs, j.volume_cu_ft, j.created_at, j.updated_at, u.username, u.status`

	var job models.Job
	var username, status string
	var avgRating float64
	
	err := r.db.QueryRow(ctx, query, jobID).Scan(
		&job.ID, &job.ContractorID, &job.ExecutorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
		&job.BulkyItems, &job.InventoryList, &job.Hoisting, &job.AdditionalServicesDescription,
		&job.EstimatedCrewAssistants, &job.TruckSize, &job.PickupAddress, &job.PickupCity, &job.PickupState, &job.PickupFloor,
		&job.PickupBuildingType, &job.PickupWalkDistance, &job.DeliveryAddress, &job.DeliveryCity, &job.DeliveryState, &job.DeliveryFloor,
		&job.DeliveryBuildingType, &job.DeliveryWalkDistance, &job.DistanceMiles, &job.JobStatus,
		&job.PickupDate, &job.PickupTimeFrom, &job.PickupTimeTo, &job.DeliveryDate,
		&job.DeliveryTimeFrom, &job.DeliveryTimeTo, &job.CutAmount, &job.PaymentAmount,
		&job.WeightLbs, &job.VolumeCuFt, &job.CreatedAt, &job.UpdatedAt,
		&username, &status, &avgRating,
	)

	if err != nil {
		return nil, err
	}

	// Assign contractor info
	job.ContractorUsername = &username
	job.ContractorStatus = &status
	job.ContractorRating = &avgRating

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
	var executorID *int64
	err = tx.QueryRow(ctx, "SELECT contractor_id, job_status, executor_id FROM jobs WHERE id = $1", jobID).Scan(&contractorID, &status, &executorID)
	if err != nil {
		return err
	}

	if contractorID == userID {
		return fmt.Errorf("you cannot claim your own job")
	}

	if status != "active" {
		return fmt.Errorf("job is not available for claiming")
	}

	if executorID != nil {
		return fmt.Errorf("job is already claimed by another user")
	}

	_, err = tx.Exec(ctx, `
		UPDATE jobs 
		SET executor_id = $1, job_status = 'claimed', updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2`,
		userID, jobID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *JobRepository) GetMyJobs(ctx context.Context, userID int64, offset, limit int) ([]models.Job, error) {
	query := `
		SELECT j.id, j.contractor_id, j.executor_id, j.job_type, j.number_of_bedrooms, j.packing_boxes, j.bulky_items,
			   j.inventory_list, j.hoisting, j.additional_services_description, j.estimated_crew_assistants,
			   j.truck_size, j.pickup_address, j.pickup_floor, j.pickup_building_type, j.pickup_walk_distance,
			   j.delivery_address, j.delivery_floor, j.delivery_building_type, j.delivery_walk_distance,
			   j.distance_miles, j.job_status, j.pickup_date, j.pickup_time_from, j.pickup_time_to,
			   j.delivery_date, j.delivery_time_from, j.delivery_time_to, j.cut_amount, j.payment_amount,
			   j.weight_lbs, j.volume_cu_ft, j.created_at, j.updated_at,
			   COALESCE(c.company_name, u.username) AS executor_name
		FROM jobs j
		LEFT JOIN users u ON j.executor_id = u.id
		LEFT JOIN companies c ON u.id = c.user_id
		WHERE j.contractor_id = $1
		ORDER BY j.created_at DESC
		LIMIT $2 OFFSET $3`

	fmt.Printf("DEBUG GetMyJobs: userID=%d, offset=%d, limit=%d\n", userID, offset, limit)
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.ContractorID, &job.ExecutorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
			&job.BulkyItems, &job.InventoryList, &job.Hoisting, &job.AdditionalServicesDescription,
			&job.EstimatedCrewAssistants, &job.TruckSize, &job.PickupAddress, &job.PickupFloor,
			&job.PickupBuildingType, &job.PickupWalkDistance, &job.DeliveryAddress, &job.DeliveryFloor,
			&job.DeliveryBuildingType, &job.DeliveryWalkDistance, &job.DistanceMiles, &job.JobStatus,
			&job.PickupDate, &job.PickupTimeFrom, &job.PickupTimeTo, &job.DeliveryDate,
			&job.DeliveryTimeFrom, &job.DeliveryTimeTo, &job.CutAmount, &job.PaymentAmount,
			&job.WeightLbs, &job.VolumeCuFt, &job.CreatedAt, &job.UpdatedAt, &job.ExecutorName,
		)
		if err != nil {
			return nil, err
		}
		fmt.Printf("DEBUG GetMyJobs: Found job ID=%d, contractor_id=%d, executor_id=%v, status=%s\n", 
			job.ID, job.ContractorID, job.ExecutorID, job.JobStatus)
		jobs = append(jobs, job)
	}

	fmt.Printf("DEBUG GetMyJobs: Total jobs found: %d\n", len(jobs))
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

func (r *JobRepository) GetPendingJobs(ctx context.Context, userID int64, limit int) ([]models.Job, error) {
	query := `
		SELECT id, contractor_id, executor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items,
			   inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			   truck_size, pickup_address, pickup_floor, pickup_building_type, pickup_walk_distance,
			   delivery_address, delivery_floor, delivery_building_type, delivery_walk_distance,
			   distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			   delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			   weight_lbs, volume_cu_ft, created_at, updated_at
		FROM jobs
		WHERE executor_id = $1 AND job_status != 'completed'
		ORDER BY pickup_date ASC, pickup_time_from ASC
		LIMIT $2`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.ContractorID, &job.ExecutorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
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

	return jobs, rows.Err()
}

func (r *JobRepository) GetCountPendingJobs(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM jobs WHERE executor_id = $1 AND job_status != 'completed'`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *JobRepository) GetClaimedJobs(ctx context.Context, userID int64, offset, limit int) ([]models.Job, error) {
	query := `
		SELECT id, contractor_id, executor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items,
			   inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			   truck_size, pickup_address, pickup_floor, pickup_building_type, pickup_walk_distance,
			   delivery_address, delivery_floor, delivery_building_type, delivery_walk_distance,
			   distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			   delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			   weight_lbs, volume_cu_ft, created_at, updated_at
		FROM jobs
		WHERE executor_id = $1
		ORDER BY updated_at DESC
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
			&job.ID, &job.ContractorID, &job.ExecutorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
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
	query := `SELECT COUNT(*) FROM jobs WHERE executor_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

// internal/repository/job_repository.go - обновить метод GetAvailableJobs

func (r *JobRepository) GetAvailableJobs(ctx context.Context, userID int64, filters *models.JobFilters) ([]models.AvailableJobDTO, int, error) {
	offset := (filters.Page - 1) * filters.Limit

	// Базовый запрос
	baseQuery := `
		SELECT id, job_type, distance_miles, pickup_address, pickup_city, pickup_state, delivery_address, delivery_city, delivery_state,
			   pickup_date, delivery_date, truck_size, weight_lbs, volume_cu_ft, payment_amount,
			   contractor_id, number_of_bedrooms, cut_amount
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL
	`

	countQuery := `
		SELECT COUNT(*) 
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL
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
		fmt.Printf("=== FILTER DEBUG: Origin filter ===\n")
		fmt.Printf("Looking for pickup location: '%s'\n", *filters.Origin)
		
		// Также добавим отладочный запрос для проверки существующих локаций
		debugPickupQuery := `SELECT pickup_city, pickup_state, pickup_city || ', ' || pickup_state as combined FROM jobs WHERE pickup_city != '' AND pickup_state != '' LIMIT 5`
		debugRows, err := r.db.Query(ctx, debugPickupQuery)
		if err == nil {
			defer debugRows.Close()
			fmt.Printf("Existing pickup data in DB:\n")
			for debugRows.Next() {
				var city, state, combined string
				if err := debugRows.Scan(&city, &state, &combined); err == nil {
					fmt.Printf("  City: '%s', State: '%s', Combined: '%s'\n", city, state, combined)
				}
			}
		}
		fmt.Printf("=== END FILTER DEBUG ===\n")
		
		conditions = append(conditions, fmt.Sprintf("pickup_city || ', ' || pickup_state = $%d", paramIndex))
		params = append(params, *filters.Origin)
		paramIndex++
	}

	if filters.Destination != nil && *filters.Destination != "" {
		conditions = append(conditions, fmt.Sprintf("delivery_city || ', ' || delivery_state = $%d", paramIndex))
		params = append(params, *filters.Destination)
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
		sizes := strings.Fields(*filters.TruckSize)
		if len(sizes) == 1 {
			conditions = append(conditions, fmt.Sprintf("truck_size = $%d", paramIndex))
			params = append(params, sizes[0])
			paramIndex++
		} else if len(sizes) > 1 {
			conditions = append(conditions, fmt.Sprintf("truck_size = ANY($%d)", paramIndex))
			params = append(params, sizes)
			paramIndex++
		}
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
	
	// Добавляем отладку для полного запроса
	fmt.Printf("=== FULL QUERY DEBUG ===\n")
	fmt.Printf("Final count query: %s\n", countQuery)
	fmt.Printf("Parameters: %v\n", params)
	
	// Проверим, есть ли записи для этого пользователя с этой локацией без других фильтров
	debugFullQuery := `
		SELECT id, contractor_id, job_status, executor_id, pickup_city, pickup_state
		FROM jobs 
		WHERE pickup_city = 'Fulshear' AND pickup_state = 'TX'
	`
	debugRows2, debugErr := r.db.Query(ctx, debugFullQuery)
	if debugErr == nil {
		defer debugRows2.Close()
		fmt.Printf("All jobs with Fulshear, TX:\n")
		for debugRows2.Next() {
			var id, contractorID int64
			var jobStatus, pickupCity, pickupState string
			var executorID *int64
			if scanErr := debugRows2.Scan(&id, &contractorID, &jobStatus, &executorID, &pickupCity, &pickupState); scanErr == nil {
				fmt.Printf("  ID: %d, ContractorID: %d, Status: %s, ExecutorID: %v\n", id, contractorID, jobStatus, executorID)
			}
		}
	}
	fmt.Printf("Current userID: %d\n", userID)
	fmt.Printf("=== END FULL QUERY DEBUG ===\n")

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
			&job.ID, &job.JobType, &job.DistanceMiles, &job.PickupAddress, &job.PickupCity, &job.PickupState,
			&job.DeliveryAddress, &job.DeliveryCity, &job.DeliveryState, &job.PickupDate, &job.DeliveryDate, &job.TruckSize,
			&job.WeightLbs, &job.VolumeCuFt, &job.PaymentAmount,
			&job.ContractorID, &job.NumberOfBedrooms, &job.CutAmount,
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
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL
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
		sizes := strings.Fields(*filters.TruckSize)
		if len(sizes) == 1 {
			conditions = append(conditions, fmt.Sprintf("truck_size = $%d", paramIndex))
			params = append(params, sizes[0])
			paramIndex++
		} else if len(sizes) > 1 {
			conditions = append(conditions, fmt.Sprintf("truck_size = ANY($%d)", paramIndex))
			params = append(params, sizes)
			paramIndex++
		}
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
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL
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
		sizes := strings.Fields(*filters.TruckSize)
		if len(sizes) == 1 {
			conditions = append(conditions, fmt.Sprintf("truck_size = $%d", paramIndex))
			params = append(params, sizes[0])
			paramIndex++
		} else if len(sizes) > 1 {
			conditions = append(conditions, fmt.Sprintf("truck_size = ANY($%d)", paramIndex))
			params = append(params, sizes)
			paramIndex++
		}
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
	options := &models.JobFilterOptions{
		NumberOfBedrooms:  []string{},
		TruckSizes:        []string{},
		PickupLocations:   []models.LocationOption{},
		DeliveryLocations: []models.LocationOption{},
	}

	// Получаем уникальные значения количества спален
	bedroomsQuery := `
		SELECT DISTINCT number_of_bedrooms 
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL 
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
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL 
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
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL
	`

	err = r.db.QueryRow(ctx, payoutRangeQuery, userID).Scan(
		&options.PayoutRange.Min,
		&options.PayoutRange.Max,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payout range: %w", err)
	}

	// Получаем уникальные pickup локации (City, State)
	pickupLocationsQuery := `
		SELECT DISTINCT pickup_city || ', ' || pickup_state as location
		FROM jobs 
		WHERE pickup_city IS NOT NULL AND pickup_city != '' 
		AND pickup_state IS NOT NULL AND pickup_state != ''
		ORDER BY location
	`

	rows, err = r.db.Query(ctx, pickupLocationsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get pickup locations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var location string
		if err := rows.Scan(&location); err != nil {
			return nil, err
		}
		options.PickupLocations = append(options.PickupLocations, models.LocationOption{
			Value: location,
			Label: location,
		})
	}

	// Получаем уникальные delivery локации (City, State)
	deliveryLocationsQuery := `
		SELECT DISTINCT delivery_city || ', ' || delivery_state as location
		FROM jobs 
		WHERE delivery_city IS NOT NULL AND delivery_city != '' 
		AND delivery_state IS NOT NULL AND delivery_state != ''
		ORDER BY location
	`

	rows, err = r.db.Query(ctx, deliveryLocationsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery locations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var location string
		if err := rows.Scan(&location); err != nil {
			return nil, err
		}
		options.DeliveryLocations = append(options.DeliveryLocations, models.LocationOption{
			Value: location,
			Label: location,
		})
	}

	// Получаем максимальную дистанцию (для слайдера), округляем до целого
	maxDistanceQuery := `
		SELECT ROUND(MAX(distance_miles))
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL
	`

	err = r.db.QueryRow(ctx, maxDistanceQuery, userID).Scan(&options.MaxDistance)
	if err != nil {
		return nil, fmt.Errorf("failed to get max distance: %w", err)
	}

	// Получаем диапазон дат
	dateRangeQuery := `
		SELECT MIN(pickup_date), MAX(pickup_date)
		FROM jobs 
		WHERE contractor_id != $1 AND job_status = 'active' AND executor_id IS NULL
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

	var currentStatus string
	var executorID *int64
	err = tx.QueryRow(ctx, "SELECT job_status, executor_id FROM jobs WHERE id = $1", jobID).Scan(&currentStatus, &executorID)
	if err != nil {
		return err
	}

	if executorID == nil || *executorID != userID {
		return fmt.Errorf("you are not the executor of this job")
	}

	if currentStatus == "completed" {
		return fmt.Errorf("job is already completed")
	}

	if currentStatus != "claimed" && currentStatus != "in_progress" && currentStatus != "pending" {
		return fmt.Errorf("job status must be claimed, in_progress, or pending to mark as completed")
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

func (r *JobRepository) CancelJobs(ctx context.Context, jobIDs []int64, userID int64) (int, error) {
	if len(jobIDs) == 0 {
		return 0, fmt.Errorf("no job IDs provided")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Check valid statuses for all jobs
	query := `
		SELECT id, job_status 
		FROM jobs 
		WHERE id = ANY($1)`

	rows, err := tx.Query(ctx, query, jobIDs)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	fmt.Printf("DEBUG: Looking for jobs with IDs: %v\n", jobIDs)
	
	validJobIDs := make([]int64, 0)
	foundJobs := make([]string, 0)
	for rows.Next() {
		var jobID int64
		var status string
		if err := rows.Scan(&jobID, &status); err != nil {
			return 0, err
		}

		foundJobs = append(foundJobs, fmt.Sprintf("ID:%d Status:%s", jobID, status))
		
		// Only allow cancelling jobs with 'active' status
		if status == "active" {
			validJobIDs = append(validJobIDs, jobID)
		}
	}
	
	fmt.Printf("DEBUG: Found jobs: %v\n", foundJobs)
	fmt.Printf("DEBUG: Jobs that can be cancelled: %v\n", validJobIDs)

	if len(validJobIDs) == 0 {
		return 0, fmt.Errorf("no jobs found with 'active' status that can be cancelled")
	}

	// Update valid jobs from 'active' to 'canceled'
	updateQuery := `
		UPDATE jobs 
		SET job_status = 'canceled', updated_at = CURRENT_TIMESTAMP 
		WHERE id = ANY($1)`

	result, err := tx.Exec(ctx, updateQuery, validJobIDs)
	if err != nil {
		return 0, err
	}

	cancelledCount := result.RowsAffected()
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return int(cancelledCount), nil
}

func (r *JobRepository) GetJobsByIDs(ctx context.Context, userID int64, jobIDs []int64) ([]models.Job, error) {
	if len(jobIDs) == 0 {
		return []models.Job{}, nil
	}

	query := `
		SELECT id, contractor_id, executor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items,
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
			&job.ID, &job.ContractorID, &job.ExecutorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
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


func (r *JobRepository) GetUserWorkStats(ctx context.Context, userID int64) (models.UserWorkStats, error) {
	var stats models.UserWorkStats

	// Получаем количество завершенных работ (где пользователь исполнитель)
	completedJobsQuery := `
		SELECT COUNT(*)
		FROM jobs
		WHERE executor_id = $1 AND job_status = 'completed'`
	
	err := r.db.QueryRow(ctx, completedJobsQuery, userID).Scan(&stats.CompletedJobs)
	if err != nil {
		return stats, fmt.Errorf("failed to get completed jobs count: %w", err)
	}

	// Получаем заработок с завершенных работ (payment_amount - cut_amount)
	earningsQuery := `
		SELECT COALESCE(SUM(payment_amount - cut_amount), 0)
		FROM jobs
		WHERE executor_id = $1 AND job_status = 'completed'`
	
	err = r.db.QueryRow(ctx, earningsQuery, userID).Scan(&stats.Earnings)
	if err != nil {
		return stats, fmt.Errorf("failed to get earnings: %w", err)
	}

	// Получаем количество предстоящих работ (claimed, in_progress)
	upcomingJobsQuery := `
		SELECT COUNT(*)
		FROM jobs
		WHERE executor_id = $1 AND job_status IN ('claimed', 'in_progress')`
	
	err = r.db.QueryRow(ctx, upcomingJobsQuery, userID).Scan(&stats.UpcomingJobs)
	if err != nil {
		return stats, fmt.Errorf("failed to get upcoming jobs count: %w", err)
	}

	return stats, nil
}

func (r *JobRepository) GetTodayScheduleJobs(ctx context.Context, userID int64, offset, limit int) ([]models.Job, error) {
	query := `
		SELECT id, contractor_id, executor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items,
			   inventory_list, hoisting, additional_services_description, estimated_crew_assistants,
			   truck_size, pickup_address, pickup_floor, pickup_building_type, pickup_walk_distance,
			   delivery_address, delivery_floor, delivery_building_type, delivery_walk_distance,
			   distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
			   delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount,
			   weight_lbs, volume_cu_ft, created_at, updated_at
		FROM jobs
		WHERE executor_id = $1 
		  AND DATE(pickup_date) = CURRENT_DATE
		ORDER BY pickup_time_from ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query today's schedule jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.ContractorID, &job.ExecutorID, &job.JobType, &job.NumberOfBedrooms, &job.PackingBoxes,
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

func (r *JobRepository) GetCountTodayScheduleJobs(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM jobs
		WHERE executor_id = $1 
		  AND DATE(pickup_date) = CURRENT_DATE`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count today's schedule jobs: %w", err)
	}

	return count, nil
}

func (r *JobRepository) InsertJobFile(ctx context.Context, jobID int64, fileID, fileName string, fileSize int64, contentType string) error {
	query := `
		INSERT INTO job_files (job_id, file_id, file_name, file_size, content_type)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, jobID, fileID, fileName, fileSize, contentType)
	return err
}

func (r *JobRepository) GetJobFiles(ctx context.Context, jobID int64) ([]models.JobFile, error) {
	query := `
		SELECT id, job_id, file_id, file_name, file_size, content_type, 
		       COALESCE(file_type, 'legacy') as file_type, uploaded_at
		FROM job_files
		WHERE job_id = $1
		ORDER BY uploaded_at DESC`

	rows, err := r.db.Query(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.JobFile
	for rows.Next() {
		var file models.JobFile
		err := rows.Scan(&file.ID, &file.JobID, &file.FileID, &file.FileName, &file.FileSize, &file.ContentType, &file.FileType, &file.UploadedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

func (r *JobRepository) UpdateJobStatus(ctx context.Context, jobID int64, status string) error {
	query := `
		UPDATE jobs 
		SET job_status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2`

	_, err := r.db.Exec(ctx, query, status, jobID)
	return err
}

func (r *JobRepository) InsertJobFileWithType(ctx context.Context, jobID int64, fileID, fileName string, fileSize int64, contentType, fileType string) error {
	query := `
		INSERT INTO job_files (job_id, file_id, file_name, file_size, content_type, file_type)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, query, jobID, fileID, fileName, fileSize, contentType, fileType)
	return err
}

func (r *JobRepository) GetJobFilesByType(ctx context.Context, jobID int64, fileType string) ([]models.JobFile, error) {
	query := `
		SELECT id, job_id, file_id, file_name, file_size, content_type, file_type, uploaded_at
		FROM job_files
		WHERE job_id = $1 AND file_type = $2
		ORDER BY uploaded_at DESC`

	rows, err := r.db.Query(ctx, query, jobID, fileType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.JobFile
	for rows.Next() {
		var file models.JobFile
		err := rows.Scan(&file.ID, &file.JobID, &file.FileID, &file.FileName, &file.FileSize, &file.ContentType, &file.FileType, &file.UploadedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}
