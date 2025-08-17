package admin

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetJobsList(ctx context.Context, limit, offset int) ([]models.JobManagementInfo, error) {
	query := `
		SELECT 
			j.id,
			j.truck_size,
			j.pickup_address,
			j.delivery_address,
			j.pickup_date,
			j.payment_amount,
			j.job_status
		FROM jobs j
		ORDER BY j.created_at DESC, j.id DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.JobManagementInfo
	for rows.Next() {
		var job models.JobManagementInfo
		err := rows.Scan(
			&job.ID,
			&job.Size,
			&job.From,
			&job.To,
			&job.Date,
			&job.Payout,
			&job.Status,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, job)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}