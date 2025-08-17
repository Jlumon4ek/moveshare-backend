package user

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) FindUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password, COALESCE(role, 'user') as role, status, profile_photo_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&user.ProfilePhotoID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
