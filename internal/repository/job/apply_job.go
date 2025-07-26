package job

import (
	"context"
	"fmt"
)

func (r *repository) ApplyJob(ctx context.Context, userID, jobID int64) error {
	query := `
		INSERT INTO job_applications (job_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT ON CONSTRAINT unique_application DO NOTHING
		RETURNING id
	`
	var id int64
	err := r.db.QueryRow(ctx, query, jobID, userID).Scan(&id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return fmt.Errorf("application already exists or job not found")
		}
		return err
	}
	return nil
}
