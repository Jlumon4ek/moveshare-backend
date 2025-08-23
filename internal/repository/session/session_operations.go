package session

import (
	"context"
)

func (r *repository) UpdateSessionActivity(ctx context.Context, sessionToken string) error {
	query := `
		UPDATE user_sessions
		SET last_activity = CURRENT_TIMESTAMP
		WHERE session_token = $1 AND expires_at > CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(ctx, query, sessionToken)
	return err
}

func (r *repository) TerminateSession(ctx context.Context, sessionID int64, userID int64) error {
	query := `
		DELETE FROM user_sessions
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.db.Exec(ctx, query, sessionID, userID)
	return err
}

func (r *repository) TerminateAllUserSessions(ctx context.Context, userID int64, exceptSessionID *int64) error {
	var query string
	var args []interface{}

	if exceptSessionID != nil {
		query = `DELETE FROM user_sessions WHERE user_id = $1 AND id != $2`
		args = []interface{}{userID, *exceptSessionID}
	} else {
		query = `DELETE FROM user_sessions WHERE user_id = $1`
		args = []interface{}{userID}
	}

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *repository) CleanupExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at < CURRENT_TIMESTAMP`

	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *repository) SetCurrentSession(ctx context.Context, sessionID int64, userID int64) error {
	// First, unset all current sessions for the user
	query1 := `
		UPDATE user_sessions
		SET is_current = FALSE
		WHERE user_id = $1
	`

	_, err := r.db.Exec(ctx, query1, userID)
	if err != nil {
		return err
	}

	// Then set the specified session as current
	query2 := `
		UPDATE user_sessions
		SET is_current = TRUE
		WHERE id = $1 AND user_id = $2
	`

	_, err = r.db.Exec(ctx, query2, sessionID, userID)
	return err
}

func (r *repository) UpdateSessionTokens(ctx context.Context, sessionID int64, accessToken, refreshToken string) error {
	query := `
		UPDATE user_sessions
		SET session_token = $1, refresh_token = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, accessToken, refreshToken, sessionID)
	return err
}