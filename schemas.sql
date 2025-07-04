CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    job_title VARCHAR(255) NOT NULL,
    description TEXT,
    cargo_type VARCHAR(50),
    urgency VARCHAR(20) NOT NULL DEFAULT 'Normal',
    truck_size VARCHAR(20),
    loading_assistance BOOLEAN DEFAULT FALSE,
    pickup_date DATE,
    pickup_time_window VARCHAR(50),
    delivery_date DATE,
    delivery_time_window VARCHAR(50),
    pickup_location VARCHAR(255) NOT NULL,
    delivery_location VARCHAR(255) NOT NULL,
    payout_amount DECIMAL(10, 2),
    early_delivery_bonus DECIMAL(10, 2),
    payment_terms VARCHAR(50),
    weight_lb DECIMAL(10, 2),
    volume_cu_ft DECIMAL(10, 2),
    liftgate BOOLEAN DEFAULT FALSE,
    fragile_items BOOLEAN DEFAULT FALSE,
    climate_control BOOLEAN DEFAULT FALSE,
    assembly_required BOOLEAN DEFAULT FALSE,
    extra_insurance BOOLEAN DEFAULT FALSE,
    additional_packing BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_applications (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT REFERENCES jobs(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_application UNIQUE (job_id, user_id)
);