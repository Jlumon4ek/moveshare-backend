package job

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUserJobs(ctx context.Context, userID int64) ([]models.Job, error) {
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
