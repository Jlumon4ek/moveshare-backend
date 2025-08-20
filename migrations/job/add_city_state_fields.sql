-- Add city and state fields to jobs table
ALTER TABLE jobs ADD COLUMN pickup_city VARCHAR(100) DEFAULT '';
ALTER TABLE jobs ADD COLUMN pickup_state VARCHAR(50) DEFAULT '';
ALTER TABLE jobs ADD COLUMN delivery_city VARCHAR(100) DEFAULT '';
ALTER TABLE jobs ADD COLUMN delivery_state VARCHAR(50) DEFAULT '';