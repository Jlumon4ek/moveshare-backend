-- Создание случайных jobs для всех пользователей
DO $$
DECLARE
    u RECORD;
    num_jobs INT;
    j INT;
    pickup_date DATE;
    delivery_date DATE;
    pickup_address_index INT;
    delivery_address_index INT;
    job_types TEXT[] := ARRAY['residential','office','warehouse','other'];
    bedrooms TEXT[] := ARRAY['1','2','3','4','5+','office'];
    crew_assistants TEXT[] := ARRAY['driver_only','driver_1','driver_2','driver_3','driver_4'];
    truck_sizes TEXT[] := ARRAY['Small','Medium','Large'];
    job_statuses TEXT[] := ARRAY['active','pending','completed','canceled'];
    addresses TEXT[] := ARRAY[
        '1600 Pennsylvania Ave NW, Washington, DC', 
        '350 5th Ave, New York, NY', 
        '1 Infinite Loop, Cupertino, CA', 
        '4059 Mt Lee Dr, Los Angeles, CA',
        '600 Montgomery St, San Francisco, CA', 
        '233 S Wacker Dr, Chicago, IL', 
        '2000 Avenue of the Stars, Los Angeles, CA', 
        '500 S Buena Vista St, Burbank, CA',
        '700 Exposition Park Dr, Los Angeles, CA', 
        '1211 Avenue of the Americas, New York, NY', 
        '1000 5th Ave, New York, NY', 
        '111 S Michigan Ave, Chicago, IL',
        '1260 6th Ave, New York, NY', 
        '100 Universal City Plaza, Universal City, CA', 
        '1 Microsoft Way, Redmond, WA', 
        '2211 N First St, San Jose, CA',
        '1601 Willow Rd, Menlo Park, CA', 
        '1 World Way, Los Angeles, CA', 
        '20 W 34th St, New York, NY', 
        '4 Yawkey Way, Boston, MA',
        '123 Main St, Fulshear, TX'
    ];
    cities TEXT[] := ARRAY[
        'Washington', 'New York', 'Cupertino', 'Los Angeles', 'San Francisco', 
        'Chicago', 'Los Angeles', 'Burbank', 'Los Angeles', 'New York',
        'New York', 'Chicago', 'New York', 'Universal City', 'Redmond',
        'San Jose', 'Menlo Park', 'Los Angeles', 'New York', 'Boston',
        'Fulshear'
    ];
    states TEXT[] := ARRAY[
        'DC', 'NY', 'CA', 'CA', 'CA', 
        'IL', 'CA', 'CA', 'CA', 'NY',
        'NY', 'IL', 'NY', 'CA', 'WA',
        'CA', 'CA', 'CA', 'NY', 'MA',
        'TX'
    ];
BEGIN
    FOR u IN SELECT * FROM users LOOP
        num_jobs := floor(random()*11 + 10)::int; -- 10–20 jobs на юзера
        FOR j IN 1..num_jobs LOOP
            pickup_date := current_date + (floor(random()*30)::int); -- сегодня + 0..29 дней
            delivery_date := pickup_date + (floor(random()*5 + 1)::int); -- delivery > pickup
            pickup_address_index := floor(random()*array_length(addresses,1)+1);
            delivery_address_index := floor(random()*array_length(addresses,1)+1);

            INSERT INTO jobs (
                contractor_id,
                job_type,
                number_of_bedrooms,
                estimated_crew_assistants,
                truck_size,
                pickup_address,
                pickup_city,
                pickup_state,
                delivery_address,
                delivery_city,
                delivery_state,
                pickup_floor,
                pickup_building_type,
                pickup_walk_distance,
                delivery_floor,
                delivery_building_type,
                delivery_walk_distance,
                distance_miles,
                job_status,
                pickup_date,
                pickup_time_from,
                pickup_time_to,
                delivery_date,
                delivery_time_from,
                delivery_time_to,
                cut_amount,
                payment_amount,
                weight_lbs,
                volume_cu_ft
            ) VALUES (
                u.id,
                job_types[floor(random()*array_length(job_types,1)+1)],
                bedrooms[floor(random()*array_length(bedrooms,1)+1)],
                crew_assistants[floor(random()*array_length(crew_assistants,1)+1)],
                truck_sizes[floor(random()*array_length(truck_sizes,1)+1)],
                addresses[pickup_address_index],
                cities[pickup_address_index],
                states[pickup_address_index],
                addresses[delivery_address_index],
                cities[delivery_address_index],
                states[delivery_address_index],
                floor(random()*20), -- pickup_floor 0..19
                'Apartment',
                (ARRAY['0-50ft','50-100ft','100-200ft'])[floor(random()*3+1)],
                floor(random()*20),
                'Apartment',
                (ARRAY['0-50ft','50-100ft','100-200ft'])[floor(random()*3+1)],
                round((random()*50 + 1)::numeric,2), -- distance_miles
                job_statuses[floor(random()*array_length(job_statuses,1)+1)],
                pickup_date,
                (time '08:00' + (random()*4 || ' hours')::interval),
                (time '12:00' + (random()*4 || ' hours')::interval),
                delivery_date,
                (time '08:00' + (random()*4 || ' hours')::interval),
                (time '12:00' + (random()*4 || ' hours')::interval),
                round((random()*200 + 50)::numeric,2), -- cut_amount
                round((random()*2000 + 500)::numeric,2), -- payment_amount
                round((random()*5000 + 100)::numeric,2), -- weight_lbs
                round((random()*500 + 50)::numeric,2) -- volume_cu_ft
            );
        END LOOP;
    END LOOP;
END $$;
