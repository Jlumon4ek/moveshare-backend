package chat

import "context"

func (r *repository) MarkMessagesAsRead(ctx context.Context, chatID, userID int64) error {
	query := `
		UPDATE chat_messages 
		SET is_read = true, read_at = NOW()
		WHERE conversation_id = $1 
		AND sender_id != $2 
		AND is_read = false
	`

	_, err := r.db.Exec(ctx, query, chatID, userID)
	return err
}
