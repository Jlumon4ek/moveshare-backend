package admin

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUsersList(ctx context.Context, limit, offset int) ([]models.UserCompanyInfo, error) {
	query := `
		SELECT 
			u.id,
			COALESCE(c.company_name, '') AS company_name,
			u.email,
			COUNT(t.id) AS trucks_number,
			u.status,
			u.created_at
		FROM users u
		LEFT JOIN companies c ON c.user_id = u.id
		LEFT JOIN trucks t ON t.user_id = u.id
		WHERE u.role = 'user'
		GROUP BY u.id, c.company_name, u.email, u.status, u.created_at
		ORDER BY u.created_at DESC, u.id DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.UserCompanyInfo
	for rows.Next() {
		var info models.UserCompanyInfo
		err := rows.Scan(
			&info.ID,
			&info.CompanyName,
			&info.Email,
			&info.TrucksNumber,
			&info.Status,
			&info.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
