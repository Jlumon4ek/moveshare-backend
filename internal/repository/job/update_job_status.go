package job

import (
	"context"
	"fmt"
)

func (r *repository) ChangeJobStatus(ctx context.Context, jobID int64, status string) error {
	query := `
		UPDATE jobs
		SET status = $1
		WHERE id = $2
	`
	result, err := r.db.Exec(ctx, query, status, jobID)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no job found with id %d", jobID)
	}
	return nil
}
