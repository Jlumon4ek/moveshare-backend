CREATE TABLE truck_photos (
    id SERIAL PRIMARY KEY,
    truck_id INTEGER NOT NULL REFERENCES trucks(id) ON DELETE CASCADE,
    photo_id TEXT NOT NULL
);
