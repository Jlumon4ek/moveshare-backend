package admin

import (
	"context"
)

func (r *repository) GetChatConversationCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM chat_conversations`
	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
