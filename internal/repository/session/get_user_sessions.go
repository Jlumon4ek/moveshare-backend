package session

import (
	"context"
	"log"
	"moveshare/internal/models"
)

func (r *repository) GetUserSessions(ctx context.Context, userID int64) ([]models.UserSession, error) {
	query := `
		SELECT id, user_id, session_token, refresh_token, user_agent, ip_address::text,
		       device_info, location_info, is_current, last_activity, expires_at,
		       created_at, updated_at
		FROM user_sessions
		WHERE user_id = $1 AND expires_at > CURRENT_TIMESTAMP
		ORDER BY last_activity DESC
	`

	log.Printf("SessionRepository: Executing query for user %d", userID)
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		log.Printf("SessionRepository: Error executing query for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var sessions []models.UserSession
	for rows.Next() {
		var session models.UserSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.SessionToken,
			&session.RefreshToken,
			&session.UserAgent,
			&session.IPAddress,
			&session.DeviceInfo,
			&session.LocationInfo,
			&session.IsCurrent,
			&session.LastActivity,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			log.Printf("SessionRepository: Error scanning row for user %d: %v", userID, err)
			return nil, err
		}
		sessions = append(sessions, session)
	}

	log.Printf("SessionRepository: Successfully retrieved %d sessions for user %d", len(sessions), userID)
	return sessions, nil
}

func (r *repository) GetSessionByToken(ctx context.Context, sessionToken string) (*models.UserSession, error) {
	query := `
		SELECT id, user_id, session_token, refresh_token, user_agent, ip_address::text,
		       device_info, location_info, is_current, last_activity, expires_at,
		       created_at, updated_at
		FROM user_sessions
		WHERE session_token = $1 AND expires_at > CURRENT_TIMESTAMP
	`

	var session models.UserSession
	err := r.db.QueryRow(ctx, query, sessionToken).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionToken,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.DeviceInfo,
		&session.LocationInfo,
		&session.IsCurrent,
		&session.LastActivity,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}