package job

import "context"

func (r *repository) JobExists(ctx context.Context, jobID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM jobs 
			WHERE id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, jobID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
