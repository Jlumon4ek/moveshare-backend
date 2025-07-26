package job

import (
	"context"
	"fmt"
)

func (r *repository) DeleteJob(ctx context.Context, userID, jobID int64) error {
	query := `
		DELETE FROM jobs
		WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(ctx, query, jobID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found or unauthorized")
	}
	return nil
}
