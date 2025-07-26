package admin

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUsersList(ctx context.Context, limit, offset int) ([]models.User, error) {
	query := `SELECT id, username, email, status, created_at FROM users WHERE role = 'user' LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
