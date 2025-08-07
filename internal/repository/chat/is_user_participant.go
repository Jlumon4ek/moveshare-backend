package chat

import "context"

// internal/repository/chat/is_user_participant.go
func (r *repository) IsUserParticipant(ctx context.Context, chatID, userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM chat_conversations cc
			WHERE cc.id = $1 
			AND (cc.client_id = $2 OR cc.contractor_id = $2)
			AND cc.status = 'active'
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
