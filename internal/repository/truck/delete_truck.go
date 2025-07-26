package truck

import (
	"context"
)

func (r *repository) DeleteTruck(ctx context.Context, id int64) error {
	query := `DELETE FROM trucks WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
