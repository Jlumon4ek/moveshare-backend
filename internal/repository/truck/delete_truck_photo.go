package truck

import "context"

func (r *repository) DeleteTruckPhoto(ctx context.Context, truckID int64, photoID string) error {
	query := `
		DELETE FROM truck_photos
		WHERE truck_id = $1 AND photo_id = $2`

	_, err := r.db.Exec(ctx, query, truckID, photoID)
	if err != nil {
		return err
	}

	return nil
}
