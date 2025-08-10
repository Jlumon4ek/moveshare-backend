package job

import "context"

func (r *repository) InsertJobPhoto(ctx context.Context, jobID int64, objectID string) error {
	query := `
		INSERT INTO job_photos (job_id, photo_id)
		VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query, jobID, objectID)
	return err
}

func (r *repository) GetJobPhotos(ctx context.Context, jobID int64) ([]string, error) {
	query := `
		SELECT photo_id
		FROM job_photos
		WHERE job_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []string
	for rows.Next() {
		var photoID string
		err := rows.Scan(&photoID)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photoID)
	}

	return photos, rows.Err()
}

func (r *repository) DeleteJobPhoto(ctx context.Context, jobID int64, photoID string) error {
	query := `
		DELETE FROM job_photos
		WHERE job_id = $1 AND photo_id = $2`

	_, err := r.db.Exec(ctx, query, jobID, photoID)
	return err
}
