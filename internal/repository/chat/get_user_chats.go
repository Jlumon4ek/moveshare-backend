package chat

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUserChats(ctx context.Context, userID int64, limit, offset int) ([]models.ChatListItem, int, error) {
	totalQuery := `
		SELECT COUNT(*)
		FROM chat_conversations cc
		WHERE cc.client_id = $1 OR cc.contractor_id = $1
		AND cc.status = 'active'
	`

	var total int
	err := r.db.QueryRow(ctx, totalQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			cc.id,
			cc.job_id,
			j.job_title,
			CASE 
				WHEN cc.client_id = $1 THEN cc.contractor_id 
				ELSE cc.client_id 
			END as other_user_id,
			CASE 
				WHEN cc.client_id = $1 THEN cu.username 
				ELSE cl.username 
			END as other_user_name,
			CASE 
				WHEN cc.client_id = $1 THEN 'contractor'
				ELSE 'client'
			END as other_user_role,
			COALESCE(cm.message_text, '') as last_message,
			COALESCE(cm.created_at, cc.created_at) as last_message_time,
			COALESCE(cm.message_type, 'text') as last_message_type,
			COALESCE(unread.unread_count, 0) as unread_count,
			COALESCE(cm.sender_id = $1, false) as is_last_msg_from_me,
			cc.status,
			cc.created_at,
			cc.updated_at
		FROM chat_conversations cc
		LEFT JOIN jobs j ON cc.job_id = j.id
		LEFT JOIN users cl ON cc.client_id = cl.id
		LEFT JOIN users cu ON cc.contractor_id = cu.id
		LEFT JOIN LATERAL (
			SELECT cm.message_text, cm.created_at, cm.message_type, cm.sender_id
			FROM chat_messages cm
			WHERE cm.conversation_id = cc.id
			ORDER BY cm.created_at DESC
			LIMIT 1
		) cm ON true
		LEFT JOIN LATERAL (
			SELECT COUNT(*) as unread_count
			FROM chat_messages cm2
			WHERE cm2.conversation_id = cc.id 
			AND cm2.sender_id != $1 
			AND cm2.is_read = false
		) unread ON true
		WHERE (cc.client_id = $1 OR cc.contractor_id = $1)
		AND cc.status = 'active'
		ORDER BY COALESCE(cm.created_at, cc.created_at) DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var chats []models.ChatListItem
	for rows.Next() {
		var chat models.ChatListItem
		err := rows.Scan(
			&chat.ID,
			&chat.JobID,
			&chat.JobTitle,
			&chat.OtherUserID,
			&chat.OtherUserName,
			&chat.OtherUserRole,
			&chat.LastMessage,
			&chat.LastMessageTime,
			&chat.LastMessageType,
			&chat.UnreadCount,
			&chat.IsLastMsgFromMe,
			&chat.Status,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}
