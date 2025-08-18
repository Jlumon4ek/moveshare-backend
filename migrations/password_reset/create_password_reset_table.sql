CREATE TABLE IF NOT EXISTS password_reset_codes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрого поиска по email и коду
CREATE INDEX idx_password_reset_email_code ON password_reset_codes(email, code);

-- Индекс для очистки истекших кодов
CREATE INDEX idx_password_reset_expires_at ON password_reset_codes(expires_at);

-- Индекс для поиска по user_id
CREATE INDEX idx_password_reset_user_id ON password_reset_codes(user_id);