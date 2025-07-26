package admin

import (
	"context"
)

func (r *repository) GetUserCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE role = 'user'`
	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
