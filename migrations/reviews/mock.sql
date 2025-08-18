INSERT INTO reviews (job_id, reviewer_id, reviewee_id, rating, comment)
SELECT
    j.id AS job_id,
    r.id AS reviewer_id,
    j.contractor_id AS reviewee_id,
    (floor(random()*5) + 1)::int AS rating,
    (ARRAY[
        'Great service, very professional!',
        'Everything went smoothly, highly recommend.',
        'The team was late but did a good job overall.',
        'Very friendly and careful with my belongings.',
        'Average experience, nothing special.',
        'Fast, efficient, and polite crew!',
        'Had some issues with communication, but resolved.',
        'Fantastic work! Will book again.',
        'Pricing was fair and service was excellent.',
        'Not satisfied with the timing, but work quality was fine.'
    ])[floor(random()*10 + 1)]
FROM jobs j
JOIN LATERAL (
    SELECT u.id
    FROM users u
    WHERE u.id <> j.contractor_id -- исключаем self-review
    ORDER BY random()
    LIMIT 1
) r ON true
WHERE NOT EXISTS (
    SELECT 1
    FROM reviews rv
    WHERE rv.job_id = j.id
      AND rv.reviewer_id = r.id
);
