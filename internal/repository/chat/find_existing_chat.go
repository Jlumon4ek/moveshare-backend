package chat

import "context"

func (r *repository) FindExistingChat(ctx context.Context, jobID, client_id, contractor_id int64) (int64, error) {
	query := `
		SELECT id 
		FROM chat_conversations 
		WHERE job_id = $1 
		AND client_id = $2
		AND contractor_id = $3
		AND status = 'active'
		LIMIT 1
	`

	var chatID int64
	err := r.db.QueryRow(ctx, query, jobID, client_id, contractor_id).Scan(&chatID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, nil
		}
		return 0, err
	}

	return chatID, nil
}
