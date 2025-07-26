package admin

import "context"

func (r *repository) GetUserRole(ctx context.Context, userID int64) (string, error) {
	query := `SELECT role FROM users WHERE id = $1`
	var role string
	err := r.db.QueryRow(ctx, query, userID).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}
