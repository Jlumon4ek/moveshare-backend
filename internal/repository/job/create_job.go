// internal/repository/job/create_job.go
package job

import (
	"context"
	"fmt"
	"moveshare/internal/models"
)

func (r *repository) CreateJob(ctx context.Context, job *models.Job, userID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO jobs (
			user_id, job_type, job_title, description, number_of_bedrooms,
			packing_boxes, bulky_items, inventory_list, hoisting, additional_services_desc,
			truck_size, crew_assistants,
			pickup_location, pickup_type, pickup_walk_distance,
			delivery_location, delivery_type, delivery_walk_distance,
			pickup_date, pickup_time_start, pickup_time_end,
			delivery_date, delivery_time_start, delivery_time_end,
			cut_amount, payment_amount, total_amount,
			status, distance_miles, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26, $27,
			$28, $29, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, query,
		userID,
		job.JobType,
		job.JobTitle,
		job.Description,
		job.NumberOfBedrooms,
		job.PackingBoxes,
		job.BulkyItems,
		job.InventoryList,
		job.Hoisting,
		job.AdditionalServicesDesc,
		job.TruckSize,
		job.CrewAssistants,
		job.PickupLocation,
		job.PickupType,
		job.PickupWalkDistance,
		job.DeliveryLocation,
		job.DeliveryType,
		job.DeliveryWalkDistance,
		job.PickupDate,
		job.PickupTimeStart,
		job.PickupTimeEnd,
		job.DeliveryDate,
		job.DeliveryTimeStart,
		job.DeliveryTimeEnd,
		job.CutAmount,
		job.PaymentAmount,
		job.TotalAmount,
		job.Status,
		job.DistanceMiles,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return tx.Commit(ctx)
}
