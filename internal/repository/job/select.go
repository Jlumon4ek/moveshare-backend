package job

import (
	"context"
	"fmt"
)

func (r *repository) GetTotalJobCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM jobs WHERE user_id != $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total job count: %w", err)
	}
	return count, nil
}
