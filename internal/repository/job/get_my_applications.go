package job

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetMyApplications(ctx context.Context, userID int64) ([]models.Job, error) {
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
