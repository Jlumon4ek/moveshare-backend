package job

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) CreateJob(ctx context.Context, job *models.Job) error {
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
