package chat

import "context"

func (r *repository) IsChatActive(ctx context.Context, chatID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM chat_conversations cc
			WHERE cc.id = $1 
			AND cc.status = 'active'
		)
	`

	var isActive bool
	err := r.db.QueryRow(ctx, query, chatID).Scan(&isActive)
	if err != nil {
		return false, err
	}

	return isActive, nil
}
