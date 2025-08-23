package session

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) CreateSession(ctx context.Context, session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions (
			user_id, session_token, refresh_token, user_agent, ip_address, 
			device_info, location_info, is_current, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at, last_activity
	`

	return r.db.QueryRow(
		ctx, query,
		session.UserID,
		session.SessionToken,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.DeviceInfo,
		session.LocationInfo,
		session.IsCurrent,
		session.ExpiresAt,
	).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.LastActivity,
	)
}