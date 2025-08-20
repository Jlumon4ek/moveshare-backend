package chat

import "context"

// GetUserUnreadCount возвращает общее количество непрочитанных сообщений для пользователя
func (r *repository) GetUserUnreadCount(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COALESCE(SUM(unread.unread_count), 0) as total_unread
		FROM chat_conversations cc
		LEFT JOIN LATERAL (
			SELECT COUNT(*) as unread_count
			FROM chat_messages cm
			WHERE cm.conversation_id = cc.id 
			AND cm.sender_id != $1 
			AND cm.is_read = false
		) unread ON true
		WHERE (cc.client_id = $1 OR cc.contractor_id = $1)
		AND cc.status = 'active'
	`

	var totalUnread int
	err := r.db.QueryRow(ctx, query, userID).Scan(&totalUnread)
	if err != nil {
		return 0, err
	}

	return totalUnread, nil
}