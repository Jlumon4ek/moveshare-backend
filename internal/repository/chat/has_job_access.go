package chat

import "context"

func (r *repository) HasJobAccess(ctx context.Context, jobID, userID1, userID2 int64) (bool, error) {
	// Получаем информацию о задании
	jobQuery := `
		SELECT contractor_id 
		FROM jobs 
		WHERE id = $1
	`

	var jobOwnerID int64
	err := r.db.QueryRow(ctx, jobQuery, jobID).Scan(&jobOwnerID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil // Задание не найдено или недоступно
		}
		return false, err
	}

	// Проверяем, что один из пользователей является владельцем задания
	if jobOwnerID == userID1 || jobOwnerID == userID2 {
		// Проверяем, что другой пользователь подавал заявку на это задание
		applicationQuery := `
			SELECT EXISTS(
				SELECT 1 
				FROM job_applications ja
				WHERE ja.job_id = $1 
				AND ja.user_id IN ($2, $3)
			)
		`

		var hasApplication bool
		err = r.db.QueryRow(ctx, applicationQuery, jobID, userID1, userID2).Scan(&hasApplication)
		if err != nil {
			return false, err
		}

		return hasApplication, nil
	}

	return false, nil
}
