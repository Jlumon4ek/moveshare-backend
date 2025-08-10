CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,

    contractor_id BIGINT REFERENCES users(id),

    job_type TEXT,
    number_of_bedrooms TEXT,

    packing_boxes BOOLEAN DEFAULT FALSE,
    bulky_items BOOLEAN DEFAULT FALSE,
    inventory_list BOOLEAN DEFAULT FALSE,
    hoisting BOOLEAN DEFAULT FALSE,
    additional_services_description TEXT,
    
    estimated_crew_assistants TEXT,
    truck_size TEXT, 
    
    pickup_address TEXT,
    pickup_floor INTEGER NULL,
    pickup_building_type TEXT,
    pickup_walk_distance TEXT,

    -- Delivery Location
    delivery_address TEXT,
    delivery_floor INTEGER NULL,
    delivery_building_type TEXT,
    delivery_walk_distance TEXT,

    -- New columns
    distance_miles DECIMAL,
    job_status TEXT,

    -- Schedule
    pickup_date DATE,
    pickup_time_from TIME,
    pickup_time_to TIME,
    delivery_date DATE,
    delivery_time_from TIME,
    delivery_time_to TIME,
    
    -- Payment Details
    cut_amount DECIMAL,
    payment_amount DECIMAL,

    -- Load details
    weight_lbs DECIMAL,
    volume_cu_ft DECIMAL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
