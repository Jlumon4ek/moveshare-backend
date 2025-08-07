package chat

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) SendMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessageResponse, error) {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Вставляем сообщение
	insertQuery := `
		INSERT INTO chat_messages (conversation_id, sender_id, message_text, message_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, insertQuery,
		message.ConversationID,
		message.SenderID,
		message.MessageText,
		message.MessageType,
	).Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Обновляем время последней активности чата
	updateChatQuery := `
		UPDATE chat_conversations 
		SET updated_at = NOW()
		WHERE id = $1
	`

	_, err = tx.Exec(ctx, updateChatQuery, message.ConversationID)
	if err != nil {
		return nil, err
	}

	// ✅ ИСПРАВЛЕНИЕ: Получаем полную информацию о сообщении для ответа
	selectQuery := `
		SELECT 
			cm.id,
			cm.sender_id,
			u.username as sender_name,
			cm.message_text,
			cm.message_type,
			cm.is_read,
			cm.read_at,
			cm.created_at,
			cm.updated_at
		FROM chat_messages cm
		LEFT JOIN users u ON cm.sender_id = u.id
		WHERE cm.id = $1
	`

	var response models.ChatMessageResponse
	err = tx.QueryRow(ctx, selectQuery, message.ID).Scan(
		&response.ID,
		&response.SenderID,
		&response.SenderName,
		&response.MessageText,
		&response.MessageType,
		&response.IsRead,
		&response.ReadAt,
		&response.CreatedAt,
		&response.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// ✅ ИСПРАВЛЕНИЕ: НЕ устанавливаем is_from_me в репозитории!
	// Это будет устанавливаться в зависимости от контекста в handler'е

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &response, nil
}
