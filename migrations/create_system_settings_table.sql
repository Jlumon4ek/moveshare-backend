-- Create system settings table
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    commission_rate DECIMAL(5,2) NOT NULL DEFAULT 7.5,
    new_user_approval VARCHAR(20) NOT NULL DEFAULT 'manual' CHECK (new_user_approval IN ('manual', 'auto')),
    minimum_payout INTEGER NOT NULL DEFAULT 500,
    job_expiration_days INTEGER NOT NULL DEFAULT 14,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default settings
INSERT INTO system_settings (commission_rate, new_user_approval, minimum_payout, job_expiration_days)
VALUES (7.5, 'manual', 500, 14)
ON CONFLICT (id) DO NOTHING;

-- Add unique constraint to ensure only one row exists
ALTER TABLE system_settings ADD CONSTRAINT single_settings_row CHECK (id = 1);