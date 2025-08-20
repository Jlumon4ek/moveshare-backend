-- Fix system_settings table
-- Remove all existing records except the one with id = 1
DELETE FROM system_settings WHERE id != 1;

-- Reset the sequence to start from 2 (so next insert gets id = 2, but we force id = 1)
SELECT setval('system_settings_id_seq', 1, false);

-- Insert default settings with id = 1 if it doesn't exist
INSERT INTO system_settings (id, commission_rate, new_user_approval, minimum_payout, job_expiration_days)
VALUES (1, 7.5, 'manual', 500, 14)
ON CONFLICT (id) DO NOTHING;