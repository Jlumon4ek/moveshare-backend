package job

import (
	"context"
	"fmt"
	"moveshare/internal/models"
)

func (r *repository) CreateJob(ctx context.Context, job *models.Job, userId int64) error {
	query := `
		INSERT INTO jobs (
			user_id, job_title, description, cargo_type, urgency,
			pickup_location, delivery_location, weight_lb, volume_cu_ft,
			truck_size, loading_assistance, pickup_date, pickup_time_window,
			delivery_date, delivery_time_window, payout_amount, early_delivery_bonus,
			payment_terms, liftgate, fragile_items, climate_control,
			assembly_required, extra_insurance, additional_packing,
			distance_miles
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25
		)`

	_, err := r.db.Exec(ctx, query,
		userId,
		job.JobTitle,
		job.Description,
		job.CargoType,
		job.Urgency,
		job.PickupLocation,
		job.DeliveryLocation,
		job.WeightLb,
		job.VolumeCuFt,
		job.TruckSize,
		job.LoadingAssistance,
		job.PickupDate,
		job.PickupTimeWindow,
		job.DeliveryDate,
		job.DeliveryTimeWindow,
		job.PayoutAmount,
		job.EarlyDeliveryBonus,
		job.PaymentTerms,
		job.Liftgate,
		job.FragileItems,
		job.ClimateControl,
		job.AssemblyRequired,
		job.ExtraInsurance,
		job.AdditionalPacking,
		job.DistanceMiles,
	)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}
