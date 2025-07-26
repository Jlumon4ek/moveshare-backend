package job

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"time"
)

func (r *repository) GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.Job, error) {
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

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	params = append(params, limit, offset)

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
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
