CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    status VARCHAR(20) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);






CREATE TABLE IF NOT EXISTS companies (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) UNIQUE,
    company_name VARCHAR(255),
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


CREATE TABLE IF NOT EXISTS trucks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    truck_name VARCHAR(255) NOT NULL,
    license_plate VARCHAR(50) NOT NULL,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    color VARCHAR(50),
    length DECIMAL(10,2),
    width DECIMAL(10,2),
    height DECIMAL(10,2),
    max_weight DECIMAL(10,2),
    truck_type VARCHAR(20) NOT NULL CHECK (truck_type IN ('Small', 'Medium', 'Large')),
    climate_control BOOLEAN DEFAULT FALSE,
    liftgate BOOLEAN DEFAULT FALSE,
    pallet_jack BOOLEAN DEFAULT FALSE,
    security_system BOOLEAN DEFAULT FALSE,
    refrigerated BOOLEAN DEFAULT FALSE,
    furniture_pads BOOLEAN DEFAULT FALSE,
    photos TEXT[], -- Array of photo URLs/paths
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_license_plate UNIQUE (license_plate)
);

-- Create index for faster queries
CREATE INDEX idx_trucks_user_id ON trucks(user_id);
CREATE INDEX idx_trucks_truck_type ON trucks(truck_type);