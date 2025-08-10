package job

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUserJobs(ctx context.Context, userID int64) ([]models.Job, error) {
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
		WHERE user_id = $1
		ORDER BY created_at DESC
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
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}
