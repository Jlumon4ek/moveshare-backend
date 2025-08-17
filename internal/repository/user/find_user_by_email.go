package user

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) FindUserByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, COALESCE(role, 'user') as role
		FROM users
		WHERE email = $1 OR username = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
