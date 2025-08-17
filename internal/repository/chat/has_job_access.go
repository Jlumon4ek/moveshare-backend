package chat

import "context"

func (r *repository) HasJobAccess(ctx context.Context, jobID, userID1, userID2 int64) (bool, error) {
	// Получаем информацию о задании (contractor_id и executor_id)
	jobQuery := `
		SELECT contractor_id, executor_id 
		FROM jobs 
		WHERE id = $1
	`

	var contractorID int64
	var executorID *int64
	err := r.db.QueryRow(ctx, jobQuery, jobID).Scan(&contractorID, &executorID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil // Задание не найдено или недоступно
		}
		return false, err
	}

	// Проверяем, что оба пользователя имеют доступ к этой работе
	// Доступ есть у: заказчика (contractor) и исполнителя (executor)
	allowedUsers := []int64{contractorID}
	if executorID != nil {
		allowedUsers = append(allowedUsers, *executorID)
	}

	// Проверяем что оба пользователя в списке разрешенных
	user1HasAccess := false
	user2HasAccess := false

	for _, allowedUserID := range allowedUsers {
		if allowedUserID == userID1 {
			user1HasAccess = true
		}
		if allowedUserID == userID2 {
			user2HasAccess = true
		}
	}

	return user1HasAccess && user2HasAccess, nil
}
