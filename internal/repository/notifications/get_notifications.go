package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"moveshare/internal/models"
	"strings"
)

func (r *repository) GetByUserID(ctx context.Context, userID int64, limit, offset int, typeFilter string, unreadOnly bool) ([]models.Notification, int, error) {
	// Build WHERE conditions
	conditions := []string{"user_id = $1"}
	args := []interface{}{userID}
	argIndex := 2

	// Add type filter
	if typeFilter != "" && typeFilter != "all" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, typeFilter)
		argIndex++
	}

	// Add unread filter
	if unreadOnly {
		conditions = append(conditions, fmt.Sprintf("is_read = $%d", argIndex))
		args = append(args, false)
		argIndex++
	}

	// Add expiration filter
	conditions = append(conditions, "(expires_at IS NULL OR expires_at > NOW())")

	whereClause := strings.Join(conditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get notifications
	query := fmt.Sprintf(`
		SELECT 
			id, user_id, type, title, message, job_id, chat_id, 
			related_user_id, is_read, priority, actions, metadata,
			created_at, read_at, expires_at
		FROM notifications 
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		var actionsJSON, metadataJSON []byte

		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message,
			&n.JobID, &n.ChatID, &n.RelatedUserID, &n.IsRead,
			&n.Priority, &actionsJSON, &metadataJSON,
			&n.CreatedAt, &n.ReadAt, &n.ExpiresAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Deserialize JSON fields
		if len(actionsJSON) > 0 {
			err = json.Unmarshal(actionsJSON, &n.Actions)
			if err != nil {
				return nil, 0, err
			}
		}

		if len(metadataJSON) > 0 {
			err = json.Unmarshal(metadataJSON, &n.Metadata)
			if err != nil {
				return nil, 0, err
			}
		}

		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *repository) GetByID(ctx context.Context, id, userID int64) (*models.Notification, error) {
	query := `
		SELECT 
			id, user_id, type, title, message, job_id, chat_id,
			related_user_id, is_read, priority, actions, metadata,
			created_at, read_at, expires_at
		FROM notifications 
		WHERE id = $1 AND user_id = $2
	`

	var n models.Notification
	var actionsJSON, metadataJSON []byte

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message,
		&n.JobID, &n.ChatID, &n.RelatedUserID, &n.IsRead,
		&n.Priority, &actionsJSON, &metadataJSON,
		&n.CreatedAt, &n.ReadAt, &n.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	// Deserialize JSON fields
	if len(actionsJSON) > 0 {
		err = json.Unmarshal(actionsJSON, &n.Actions)
		if err != nil {
			return nil, err
		}
	}

	if len(metadataJSON) > 0 {
		err = json.Unmarshal(metadataJSON, &n.Metadata)
		if err != nil {
			return nil, err
		}
	}

	return &n, nil
}