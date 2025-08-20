package chat

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetChatParticipants(ctx context.Context, chatID int64) ([]models.ChatParticipant, error) {
	query := `
		SELECT 
			cc.client_id as user_id,
			uc.username as user_name,
			'client' as role
		FROM chat_conversations cc
		JOIN users uc ON uc.id = cc.client_id
		WHERE cc.id = $1
		
		UNION ALL
		
		SELECT 
			cc.contractor_id as user_id,
			uco.username as user_name,
			'contractor' as role
		FROM chat_conversations cc
		JOIN users uco ON uco.id = cc.contractor_id
		WHERE cc.id = $1
	`

	rows, err := r.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.ChatParticipant
	for rows.Next() {
		var participant models.ChatParticipant
		err := rows.Scan(
			&participant.UserID,
			&participant.UserName,
			&participant.Role,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, rows.Err()
}