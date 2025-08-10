package chat

import (
	"context"
	"errors"
)

func (r *repository) CreateChat(ctx context.Context, jobID, clientID, contractorID int64) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Определяем роли пользователей
	var actualClientID, actualContractorID int64

	// Получаем информацию о задании, чтобы понять кто клиент
	jobOwnerQuery := `SELECT user_id FROM jobs WHERE id = $1`
	var jobOwnerID int64
	err = tx.QueryRow(ctx, jobOwnerQuery, jobID).Scan(&jobOwnerID)
	if err != nil {
		return 0, err
	}

	// Определяем кто клиент, а кто подрядчик
	switch jobOwnerID {
	case clientID:
		actualClientID = clientID
		actualContractorID = contractorID
	case contractorID:
		actualClientID = contractorID
		actualContractorID = clientID
	default:
		return 0, errors.New("one of the participants must be the job owner")
	}

	// Создаем чат
	insertQuery := `
		INSERT INTO chat_conversations (job_id, client_id, contractor_id, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id
	`

	var chatID int64
	err = tx.QueryRow(ctx, insertQuery, jobID, actualClientID, actualContractorID).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return chatID, nil
}
