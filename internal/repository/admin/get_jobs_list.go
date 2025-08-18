package admin

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"strings"
)

func (r *repository) GetJobsList(ctx context.Context, limit, offset int, statuses []string) ([]models.JobManagementInfo, error) {
	var query string
	var args []interface{}
	
	if len(statuses) == 0 {
		query = `
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
		args = []interface{}{limit, offset}
	} else {
		placeholders := make([]string, len(statuses))
		args = make([]interface{}, len(statuses)+2)
		
		for i, status := range statuses {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args[i] = status
		}
		args[len(statuses)] = limit
		args[len(statuses)+1] = offset
		
		query = fmt.Sprintf(`
			SELECT 
				j.id,
				j.truck_size,
				j.pickup_address,
				j.delivery_address,
				j.pickup_date,
				j.payment_amount,
				j.job_status
			FROM jobs j
			WHERE j.job_status IN (%s)
			ORDER BY j.created_at DESC, j.id DESC
			LIMIT $%d OFFSET $%d;
		`, strings.Join(placeholders, ","), len(statuses)+1, len(statuses)+2)
	}

	rows, err := r.db.Query(ctx, query, args...)
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