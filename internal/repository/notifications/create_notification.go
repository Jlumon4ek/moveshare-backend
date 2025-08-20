package notifications

import (
	"context"
	"encoding/json"
	"moveshare/internal/models"
)

func (r *repository) Create(ctx context.Context, req *models.NotificationRequest) (*models.Notification, error) {
	// Serialize actions and metadata to JSON
	actionsJSON, err := json.Marshal(req.Actions)
	if err != nil {
		return nil, err
	}

	metadataJSON, err := json.Marshal(req.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO notifications (
			user_id, type, title, message, job_id, chat_id, 
			related_user_id, priority, actions, metadata, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, created_at
	`

	var notification models.Notification
	err = r.db.QueryRow(ctx, query,
		req.UserID, req.Type, req.Title, req.Message, req.JobID,
		req.ChatID, req.RelatedUserID, req.Priority, actionsJSON,
		metadataJSON, req.ExpiresAt,
	).Scan(&notification.ID, &notification.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Fill in the rest of the fields
	notification.UserID = req.UserID
	notification.Type = req.Type
	notification.Title = req.Title
	notification.Message = req.Message
	notification.JobID = req.JobID
	notification.ChatID = req.ChatID
	notification.RelatedUserID = req.RelatedUserID
	notification.Priority = req.Priority
	notification.Actions = req.Actions
	notification.Metadata = req.Metadata
	notification.IsRead = false
	notification.ExpiresAt = req.ExpiresAt

	return &notification, nil
}