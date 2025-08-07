package chat

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetChatMessages(ctx context.Context, chatID, userID int64, limit, offset int, order string) ([]models.ChatMessageResponse, int, error) {
	// Сначала получаем общее количество сообщений
	totalQuery := `
		SELECT COUNT(*)
		FROM chat_messages cm
		WHERE cm.conversation_id = $1
	`

	var total int
	err := r.db.QueryRow(ctx, totalQuery, chatID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Определяем направление сортировки
	orderClause := "ORDER BY cm.created_at DESC"
	if order == "asc" {
		orderClause = "ORDER BY cm.created_at ASC"
	}

	// Основной запрос для получения сообщений
	query := `
		SELECT 
			cm.id,
			cm.sender_id,
			u.username as sender_name,
			cm.message_text,
			cm.message_type,
			cm.is_read,
			cm.read_at,
			(cm.sender_id = $2) as is_from_me,
			cm.created_at,
			cm.updated_at
		FROM chat_messages cm
		LEFT JOIN users u ON cm.sender_id = u.id
		WHERE cm.conversation_id = $1
		` + orderClause + `
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, chatID, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []models.ChatMessageResponse
	for rows.Next() {
		var msg models.ChatMessageResponse
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.SenderName,
			&msg.MessageText,
			&msg.MessageType,
			&msg.IsRead,
			&msg.ReadAt,
			&msg.IsFromMe,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}
