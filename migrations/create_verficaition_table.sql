CREATE TABLE file_ids (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    object_name TEXT NOT NULL,
    file_type TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_user_filetype
        UNIQUE(user_id, file_type)
);