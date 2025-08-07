package chat

import "context"

func (r *repository) UpdateChatActivity(ctx context.Context, chatID int64) error {
	query := `
		UPDATE chat_conversations 
		SET updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, chatID)
	return err
}
