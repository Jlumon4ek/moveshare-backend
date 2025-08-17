package admin

import (
	"context"
)

func (r *repository) GetActiveJobsCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM jobs WHERE job_status != 'completed'`
	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}