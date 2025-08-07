package job

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"strconv"
	"time"
)

func (r *repository) GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.Job, error) {
	query := `
		SELECT id, user_id, job_title, description, cargo_type, urgency, 
		       pickup_location, delivery_location, weight_lb, volume_cu_ft, 
		       truck_size, loading_assistance, pickup_date, pickup_time_window, 
		       delivery_date, delivery_time_window, payout_amount, early_delivery_bonus, 
		       payment_terms, liftgate, fragile_items, climate_control, 
		       assembly_required, extra_insurance, additional_packing, 
		       status, created_at, updated_at, distance_miles
		FROM jobs 
		WHERE user_id != $1`

	args := []interface{}{userID}
	argIndex := 2

	if origin, exists := filters["origin"]; exists && origin != "" {
		query += ` AND pickup_location ILIKE $` + fmt.Sprintf("%d", argIndex)
		args = append(args, "%"+origin+"%")
		argIndex++
	}

	if destination, exists := filters["destination"]; exists && destination != "" {
		query += ` AND delivery_location ILIKE $` + fmt.Sprintf("%d", argIndex)
		args = append(args, "%"+destination+"%")
		argIndex++
	}

	if distance, exists := filters["distance"]; exists && distance != "" {
		if distanceValue, err := strconv.ParseFloat(distance, 64); err == nil {
			query += ` AND distance_miles <= $` + fmt.Sprintf("%d", argIndex)
			args = append(args, distanceValue)
			argIndex++
		}
	}

	if dateStart, exists := filters["date_start"]; exists && dateStart != "" {
		if startDate, err := time.Parse("01/02/2006", dateStart); err == nil {
			query += ` AND pickup_date >= $` + fmt.Sprintf("%d", argIndex)
			args = append(args, startDate)
			argIndex++
		}
	}

	if dateEnd, exists := filters["date_end"]; exists && dateEnd != "" {
		if endDate, err := time.Parse("01/02/2006", dateEnd); err == nil {
			query += ` AND pickup_date <= $` + fmt.Sprintf("%d", argIndex)
			args = append(args, endDate)
			argIndex++
		}
	}

	if truckSize, exists := filters["truck_size"]; exists && truckSize != "" {
		query += ` AND truck_size = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, truckSize)
		argIndex++
	}

	if minPayout, exists := filters["min_payout"]; exists && minPayout != "" {
		if minValue, err := strconv.ParseFloat(minPayout, 64); err == nil {
			query += ` AND payout_amount >= $` + fmt.Sprintf("%d", argIndex)
			args = append(args, minValue)
			argIndex++
		}
	}

	if maxPayout, exists := filters["max_payout"]; exists && maxPayout != "" {
		if maxValue, err := strconv.ParseFloat(maxPayout, 64); err == nil {
			query += ` AND payout_amount <= $` + fmt.Sprintf("%d", argIndex)
			args = append(args, maxValue)
			argIndex++
		}
	}

	query += ` ORDER BY created_at DESC`

	query += ` LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query available jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.UserID, &job.JobTitle, &job.Description,
			&job.CargoType, &job.Urgency, &job.PickupLocation,
			&job.DeliveryLocation, &job.WeightLb, &job.VolumeCuFt,
			&job.TruckSize, &job.LoadingAssistance, &job.PickupDate,
			&job.PickupTimeWindow, &job.DeliveryDate, &job.DeliveryTimeWindow,
			&job.PayoutAmount, &job.EarlyDeliveryBonus, &job.PaymentTerms,
			&job.Liftgate, &job.FragileItems, &job.ClimateControl,
			&job.AssemblyRequired, &job.ExtraInsurance, &job.AdditionalPacking,
			&job.Status, &job.CreatedAt, &job.UpdatedAt, &job.DistanceMiles,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job row: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over job rows: %w", err)
	}

	return jobs, nil
}
