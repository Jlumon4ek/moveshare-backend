package admin

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetSystemSettings(ctx context.Context) (*models.SystemSettings, error) {
	// First, try to create the table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS system_settings (
			id SERIAL PRIMARY KEY,
			commission_rate DECIMAL(5,2) NOT NULL DEFAULT 7.5,
			new_user_approval VARCHAR(20) NOT NULL DEFAULT 'manual' CHECK (new_user_approval IN ('manual', 'auto')),
			minimum_payout INTEGER NOT NULL DEFAULT 500,
			job_expiration_days INTEGER NOT NULL DEFAULT 14,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	
	_, err := r.db.Exec(ctx, createTableQuery)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, commission_rate, new_user_approval, minimum_payout, job_expiration_days
		FROM system_settings 
		WHERE id = 1
	`

	var settings models.SystemSettings
	err = r.db.QueryRow(ctx, query).Scan(
		&settings.ID,
		&settings.CommissionRate,
		&settings.NewUserApproval,
		&settings.MinimumPayout,
		&settings.JobExpirationDays,
	)

	if err != nil {
		// If no settings exist, return default values
		return &models.SystemSettings{
			ID:                1,
			CommissionRate:    7.5,
			NewUserApproval:   "manual",
			MinimumPayout:     500,
			JobExpirationDays: 14,
		}, nil
	}

	return &settings, nil
}

func (r *repository) UpdateSystemSettings(ctx context.Context, settings *models.SystemSettings) error {
	// Always update the first (and should be only) record
	// Use UPSERT to either insert or update
	query := `
		INSERT INTO system_settings (id, commission_rate, new_user_approval, minimum_payout, job_expiration_days)
		VALUES (1, $1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			commission_rate = EXCLUDED.commission_rate,
			new_user_approval = EXCLUDED.new_user_approval,
			minimum_payout = EXCLUDED.minimum_payout,
			job_expiration_days = EXCLUDED.job_expiration_days,
			updated_at = NOW()
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		settings.CommissionRate,
		settings.NewUserApproval,
		settings.MinimumPayout,
		settings.JobExpirationDays,
	).Scan(&settings.ID)

	return err
}