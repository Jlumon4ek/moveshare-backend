package user

import (
	"context"
	"fmt"
)

func (r *repository) UpdateProfilePhotoID(ctx context.Context, userID int64, photoID string) error {
	var query string
	var args []interface{}

	if photoID == "" {
		// Set profile_photo_id to NULL
		query = `
			UPDATE users 
			SET profile_photo_id = NULL, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $1
		`
		args = []interface{}{userID}
	} else {
		// Set profile_photo_id to photoID
		query = `
			UPDATE users 
			SET profile_photo_id = $1, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $2
		`
		args = []interface{}{photoID, userID}
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile photo ID: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}