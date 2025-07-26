package job

import (
	"context"
	"fmt"
)

func (r *repository) ApplyJob(ctx context.Context, userID int64, jobID int64) error {
	query := `
		INSERT INTO job_applications (job_id, user_id)
		VALUES ($1, $2)
	`
	_, err := r.db.Exec(ctx, query, jobID, userID)
	if err != nil {
		return fmt.Errorf("failed to apply for job: %w", err)
	}
	return nil
}
