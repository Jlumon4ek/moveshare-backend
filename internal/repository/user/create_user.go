package user

import (
	"context"
	"moveshare/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (username, email, password, role)
		VALUES ($1, $2, $3, COALESCE($4, 'user'))
		RETURNING id
	`

	role := user.Role
	if role == "" {
		role = "user"
	}

	return r.db.QueryRow(ctx, query, user.Username, user.Email, string(hashedPassword), role).Scan(&user.ID)
}
