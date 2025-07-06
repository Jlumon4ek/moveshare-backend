-- Add unique constraint to trucks license_plate field
-- This migration ensures license plates are unique across all trucks

ALTER TABLE trucks ADD CONSTRAINT trucks_license_plate_unique UNIQUE (license_plate);