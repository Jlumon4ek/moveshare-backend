CREATE TABLE trucks (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    truck_name TEXT NOT NULL,
    license_plate TEXT NOT NULL,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    year INT NOT NULL,
    color TEXT NOT NULL,
    length DOUBLE PRECISION NOT NULL,
    width DOUBLE PRECISION NOT NULL,
    height DOUBLE PRECISION NOT NULL,
    max_weight DOUBLE PRECISION NOT NULL,
    truck_type TEXT NOT NULL CHECK (truck_type IN ('Small', 'Medium', 'Large')),
    climate_control BOOLEAN DEFAULT FALSE,
    liftgate BOOLEAN DEFAULT FALSE,
    pallet_jack BOOLEAN DEFAULT FALSE,
    security_system BOOLEAN DEFAULT FALSE,
    refrigerated BOOLEAN DEFAULT FALSE,
    furniture_pads BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

);
