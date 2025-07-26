package truck

import "context"

func (r *repository) GetTruckPhotos(ctx context.Context, truckID int64) ([]string, error) {
	query := `
		SELECT photo_id
		FROM truck_photos
		WHERE truck_id = $1`

	rows, err := r.db.Query(ctx, query, truckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []string
	for rows.Next() {
		var photoURL string
		err := rows.Scan(&photoURL)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photoURL)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}
