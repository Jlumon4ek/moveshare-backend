INSERT INTO companies (
    user_id,
    company_name,
    email_address,
    address,
    state,
    mc_license_number,
    company_description,
    contact_person,
    phone_number,
    city,
    zip_code,
    dot_number
)
SELECT
    u.id,
    'Company ' || u.username,
    lower(u.username) || '@businessmail.com',
    (ARRAY[
        '1600 Pennsylvania Ave NW',
        '350 5th Ave',
        '1 Infinite Loop',
        '4059 Mt Lee Dr',
        '600 Montgomery St',
        '233 S Wacker Dr',
        '2000 Avenue of the Stars',
        '500 S Buena Vista St',
        '700 Exposition Park Dr',
        '1211 Avenue of the Americas',
        '1000 5th Ave',
        '111 S Michigan Ave',
        '1260 6th Ave',
        '100 Universal City Plaza',
        '1 Microsoft Way',
        '2211 N First St',
        '1601 Willow Rd',
        '1 World Way',
        '20 W 34th St',
        '4 Yawkey Way'
    ])[floor(random()*20 + 1)],
    (ARRAY['CA','NY','IL','TX','FL','MA','WA','NJ','DC','NV'])[floor(random()*10 + 1)],
    'MC' || floor(random()*900000 + 100000), -- MC license number
    'Reliable moving services for residential and office clients.',
    initcap(split_part(u.username, '.', 1)) || ' ' || initcap(split_part(u.username, '.', 2)),
    '+1' || (floor(random()*9000000000 + 1000000000))::text,
    (ARRAY['Los Angeles','New York','Chicago','San Francisco','Boston','Miami','Seattle','Houston','Dallas','Atlanta'])[floor(random()*10 + 1)],
    lpad(floor(random()*90000 + 10000)::text, 5, '0'),
    'DOT' || floor(random()*900000 + 100000)
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM companies c WHERE c.user_id = u.id
);
