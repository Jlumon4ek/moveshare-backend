package verification

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) SelectVerificationFiles(ctx context.Context, userID int64) ([]models.VerificationFile, error) {
	query := `
		SELECT object_name, file_type, status
		FROM verification_file
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.VerificationFile
	for rows.Next() {
		var file models.VerificationFile
		if err := rows.Scan(&file.ObjectName, &file.FileType, &file.Status); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
