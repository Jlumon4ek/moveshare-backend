package job

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetMyApplications(ctx context.Context, userID int64) ([]models.Job, error) {
	query := `
		SELECT 
			j.id, j.user_id, j.job_type, j.job_title, j.description, j.number_of_bedrooms,
			j.packing_boxes, j.bulky_items, j.inventory_list, j.hoisting, j.additional_services_desc,
			j.truck_size, j.crew_assistants,
			j.pickup_location, j.pickup_type, j.pickup_walk_distance,
			j.delivery_location, j.delivery_type, j.delivery_walk_distance,
			j.pickup_date, j.pickup_time_start, j.pickup_time_end,
			j.delivery_date, j.delivery_time_start, j.delivery_time_end,
			j.cut_amount, j.payment_amount, j.total_amount,
			j.status, j.distance_miles, j.created_at, j.updated_at
		FROM jobs j
		JOIN job_applications ja ON j.id = ja.job_id
		WHERE ja.user_id = $1
		ORDER BY ja.created_at DESC
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
