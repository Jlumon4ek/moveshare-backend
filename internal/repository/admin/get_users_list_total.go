package admin

import (
	"context"
)

func (r *repository) GetUsersListTotal(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		LEFT JOIN companies c ON c.user_id = u.id
		LEFT JOIN trucks t ON t.user_id = u.id
		WHERE u.role = 'user';
	`

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}