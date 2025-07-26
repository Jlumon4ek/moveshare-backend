package verification

import (
	"context"
	"fmt"
)

func (r *repository) InsertFileID(ctx context.Context, userID int64, objectName string, fileType string) error {
	query := `
		INSERT INTO verification_file (user_id, object_name, file_type)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, file_type) DO UPDATE
		SET object_name = EXCLUDED.object_name
	`

	_, err := r.db.Exec(ctx, query, userID, objectName, fileType)
	if err != nil {
		return fmt.Errorf("failed to insert file ID: %w", err)
	}
	return nil
}
