package notifications

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) MarkAsRead(ctx context.Context, id, userID int64) error {
	query := `
		UPDATE notifications 
		SET is_read = true, read_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}

func (r *repository) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := `
		UPDATE notifications 
		SET is_read = true, read_at = NOW()
		WHERE user_id = $1 AND is_read = false
	`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *repository) Delete(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM notifications WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}

func (r *repository) DeleteAll(ctx context.Context, userID int64) error {
	query := `DELETE FROM notifications WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *repository) GetStats(ctx context.Context, userID int64) (*models.NotificationStats, error) {
	// Get total and unread counts
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN is_read = false THEN 1 END) as unread
		FROM notifications 
		WHERE user_id = $1 AND (expires_at IS NULL OR expires_at > NOW())
	`

	var stats models.NotificationStats
	err := r.db.QueryRow(ctx, query, userID).Scan(&stats.Total, &stats.Unread)
	if err != nil {
		return nil, err
	}

	// Get counts by type
	typeQuery := `
		SELECT type, COUNT(*) 
		FROM notifications 
		WHERE user_id = $1 AND (expires_at IS NULL OR expires_at > NOW())
		GROUP BY type
	`

	rows, err := r.db.Query(ctx, typeQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.ByType = make(map[models.NotificationType]int)
	for rows.Next() {
		var notType models.NotificationType
		var count int
		err := rows.Scan(&notType, &count)
		if err != nil {
			return nil, err
		}
		stats.ByType[notType] = count
	}

	return &stats, nil
}

func (r *repository) CleanupExpired(ctx context.Context) (int, error) {
	query := `DELETE FROM notifications WHERE expires_at < NOW()`
	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}