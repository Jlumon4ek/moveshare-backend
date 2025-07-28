package user

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) FindUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
