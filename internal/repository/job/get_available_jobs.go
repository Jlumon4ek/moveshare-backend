package job

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"time"
)

func (r *repository) GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.Job, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM jobs 
		WHERE user_id != $1 AND status = 'active'
	`

	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with filters
	query := `
		SELECT 
			id, user_id, job_type, job_title, description, number_of_bedrooms,
			packing_boxes, bulky_items, inventory_list, hoisting, additional_services_desc,
			truck_size, crew_assistants,
			pickup_location, pickup_type, pickup_walk_distance,
			delivery_location, delivery_type, delivery_walk_distance,
			pickup_date, pickup_time_start, pickup_time_end,
			delivery_date, delivery_time_start, delivery_time_end,
			cut_amount, payment_amount, total_amount,
			status, distance_miles, created_at, updated_at
		FROM jobs 
		WHERE user_id != $1 AND status = 'active'
	`

	args := []interface{}{userID}
	argIndex := 2

	// Apply filters
	if jobType, exists := filters["job_type"]; exists && jobType != "" {
		query += fmt.Sprintf(` AND job_type = $%d`, argIndex)
		args = append(args, jobType)
		argIndex++
	}

	if pickupLocation, exists := filters["pickup_location"]; exists && pickupLocation != "" {
		query += fmt.Sprintf(` AND pickup_location ILIKE $%d`, argIndex)
		args = append(args, "%"+pickupLocation+"%")
		argIndex++
	}

	if deliveryLocation, exists := filters["delivery_location"]; exists && deliveryLocation != "" {
		query += fmt.Sprintf(` AND delivery_location ILIKE $%d`, argIndex)
		args = append(args, "%"+deliveryLocation+"%")
		argIndex++
	}

	if truckSize, exists := filters["truck_size"]; exists && truckSize != "" {
		query += fmt.Sprintf(` AND truck_size = $%d`, argIndex)
		args = append(args, truckSize)
		argIndex++
	}

	if bedrooms, exists := filters["number_of_bedrooms"]; exists && bedrooms != "" {
		query += fmt.Sprintf(` AND number_of_bedrooms = $%d`, argIndex)
		args = append(args, bedrooms)
		argIndex++
	}

	// Date filters
	if pickupDateStart, exists := filters["pickup_date_start"]; exists && pickupDateStart != "" {
		if startDate, err := time.Parse("2006-01-02", pickupDateStart); err == nil {
			query += fmt.Sprintf(` AND pickup_date >= $%d`, argIndex)
			args = append(args, startDate)
			argIndex++
		}
	}

	if pickupDateEnd, exists := filters["pickup_date_end"]; exists && pickupDateEnd != "" {
		if endDate, err := time.Parse("2006-01-02", pickupDateEnd); err == nil {
			query += fmt.Sprintf(` AND pickup_date <= $%d`, argIndex)
			args = append(args, endDate)
			argIndex++
		}
	}

	// Payment filters
	if minPayment, exists := filters["min_payment"]; exists && minPayment != "" {
		query += fmt.Sprintf(` AND payment_amount >= $%d`, argIndex)
		args = append(args, minPayment)
		argIndex++
	}

	if maxPayment, exists := filters["max_payment"]; exists && maxPayment != "" {
		query += fmt.Sprintf(` AND payment_amount <= $%d`, argIndex)
		args = append(args, maxPayment)
		argIndex++
	}

	query += ` ORDER BY created_at DESC`
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query available jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.UserID, &job.JobType, &job.JobTitle, &job.Description, &job.NumberOfBedrooms,
			&job.PackingBoxes, &job.BulkyItems, &job.InventoryList, &job.Hoisting, &job.AdditionalServicesDesc,
			&job.TruckSize, &job.CrewAssistants,
			&job.PickupLocation, &job.PickupType, &job.PickupWalkDistance,
			&job.DeliveryLocation, &job.DeliveryType, &job.DeliveryWalkDistance,
			&job.PickupDate, &job.PickupTimeStart, &job.PickupTimeEnd,
			&job.DeliveryDate, &job.DeliveryTimeStart, &job.DeliveryTimeEnd,
			&job.CutAmount, &job.PaymentAmount, &job.TotalAmount,
			&job.Status, &job.DistanceMiles, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan job row: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over job rows: %w", err)
	}

	return jobs, total, nil
}
