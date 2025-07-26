package truck

import (
	"context"
)

func (r *repository) InsertPhoto(ctx context.Context, truckID int64, objectID string) error {
	query := `
		INSERT INTO truck_photos (truck_id, photo_id)
		VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query, truckID, objectID)
	if err != nil {
		return err
	}

	return nil
}
