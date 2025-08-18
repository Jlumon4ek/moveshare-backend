INSERT INTO jobs (
    contractor_id,
    job_type,
    number_of_bedrooms,
    packing_boxes,
    bulky_items,
    inventory_list,
    hoisting,
    additional_services_description,
    estimated_crew_assistants,
    truck_size,
    pickup_address,
    pickup_floor,
    pickup_building_type,
    pickup_walk_distance,
    delivery_address,
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
)
SELECT
    u.id AS contractor_id,
    -- job_type
    (ARRAY['residential','office','warehouse','other'])[floor(random() * 4 + 1)],
    -- number_of_bedrooms
    (ARRAY['1','2','3','4','5+','office'])[floor(random() * 6 + 1)],
    -- boolean flags
    (random() < 0.5),
    (random() < 0.5),
    (random() < 0.5),
    (random() < 0.5),
    'Additional service details ' || floor(random()*100),
    -- estimated_crew_assistants
    (ARRAY['driver_only','driver_1','driver_2','driver_3','driver_4'])[floor(random() * 5 + 1)],
    -- truck_size
    (ARRAY['Small','Medium','Large'])[floor(random() * 3 + 1)],

    -- pickup_address (реальные адреса из США)
    (ARRAY[
        '1600 Pennsylvania Ave NW, Washington, DC 20500',
        '350 5th Ave, New York, NY 10118',
        '1 Infinite Loop, Cupertino, CA 95014',
        '4059 Mt Lee Dr, Hollywood, CA 90068',
        '600 Montgomery St, San Francisco, CA 94111',
        '233 S Wacker Dr, Chicago, IL 60606',
        '2000 Avenue of the Stars, Los Angeles, CA 90067',
        '500 S Buena Vista St, Burbank, CA 91521',
        '700 Exposition Park Dr, Los Angeles, CA 90037',
        '1211 Avenue of the Americas, New York, NY 10036',
        '1000 5th Ave, New York, NY 10028',
        '111 S Michigan Ave, Chicago, IL 60603',
        '1260 6th Ave, New York, NY 10020',
        '100 Universal City Plaza, Universal City, CA 91608',
        '1 Microsoft Way, Redmond, WA 98052',
        '2211 N First St, San Jose, CA 95131',
        '1601 Willow Rd, Menlo Park, CA 94025',
        '1 World Way, Los Angeles, CA 90045',
        '20 W 34th St, New York, NY 10001',
        '4 Yawkey Way, Boston, MA 02215'
    ])[floor(random()*20 + 1)],

    floor(random()*10),
    (ARRAY['Apartment','House','Office','Warehouse'])[floor(random()*4 + 1)],
    (ARRAY['Short','Medium','Long'])[floor(random()*3 + 1)],

    -- delivery_address (тоже реальные из списка)
    (ARRAY[
        '1600 Pennsylvania Ave NW, Washington, DC 20500',
        '350 5th Ave, New York, NY 10118',
        '1 Infinite Loop, Cupertino, CA 95014',
        '4059 Mt Lee Dr, Hollywood, CA 90068',
        '600 Montgomery St, San Francisco, CA 94111',
        '233 S Wacker Dr, Chicago, IL 60606',
        '2000 Avenue of the Stars, Los Angeles, CA 90067',
        '500 S Buena Vista St, Burbank, CA 91521',
        '700 Exposition Park Dr, Los Angeles, CA 90037',
        '1211 Avenue of the Americas, New York, NY 10036',
        '1000 5th Ave, New York, NY 10028',
        '111 S Michigan Ave, Chicago, IL 60603',
        '1260 6th Ave, New York, NY 10020',
        '100 Universal City Plaza, Universal City, CA 91608',
        '1 Microsoft Way, Redmond, WA 98052',
        '2211 N First St, San Jose, CA 95131',
        '1601 Willow Rd, Menlo Park, CA 94025',
        '1 World Way, Los Angeles, CA 90045',
        '20 W 34th St, New York, NY 10001',
        '4 Yawkey Way, Boston, MA 02215'
    ])[floor(random()*20 + 1)],

    floor(random()*10),
    (ARRAY['Apartment','House','Office','Warehouse'])[floor(random()*4 + 1)],
    (ARRAY['Short','Medium','Long'])[floor(random()*3 + 1)],

    round((random() * 50)::numeric, 2),
    (ARRAY['active','completed','pending','canceled'])[floor(random() * 4 + 1)],

    -- schedule
    CURRENT_DATE + (trunc(random()*30)::int),
    time '08:00' + (trunc(random()*8) * interval '1 hour'),
    time '12:00' + (trunc(random()*8) * interval '1 hour'),
    CURRENT_DATE + (trunc(random()*30)::int),
    time '13:00' + (trunc(random()*8) * interval '1 hour'),
    time '18:00' + (trunc(random()*8) * interval '1 hour'),

    -- payments
    round((random() * 200)::numeric, 2),
    round((random() * 1000 + 200)::numeric, 2),

    -- load
    round((random() * 5000)::numeric, 2),
    round((random() * 800)::numeric, 2)
FROM users u
JOIN generate_series(1, (10 + floor(random()*11))::int) g ON true
WHERE u.id BETWEEN 9 AND 28;
