CREATE TABLE IF NOT EXISTS companies (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) UNIQUE,
    email_address VARCHAR(255),
    address TEXT,
    state VARCHAR(50),
    mc_license_number VARCHAR(20),
    company_description TEXT,
    contact_person VARCHAR(100),
    phone_number VARCHAR(20),
    city VARCHAR(100),
    zip_code VARCHAR(10),
    dot_number VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
