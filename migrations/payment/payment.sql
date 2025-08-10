-- migrations/add_stripe_tables.sql

-- Добавляем stripe_customer_id к пользователям
ALTER TABLE users ADD COLUMN stripe_customer_id VARCHAR(255) UNIQUE;

-- Таблица для сохраненных карт пользователей
CREATE TABLE user_payment_methods (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    stripe_payment_method_id VARCHAR(255) NOT NULL,
    stripe_customer_id VARCHAR(255) NOT NULL,
    card_last4 VARCHAR(4) NOT NULL,
    card_brand VARCHAR(20) NOT NULL,
    card_exp_month INT NOT NULL,
    card_exp_year INT NOT NULL,
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(user_id, stripe_payment_method_id)
);

-- Уникальный индекс только для активной дефолтной карты
CREATE UNIQUE INDEX uniq_user_default_card
ON user_payment_methods(user_id)
WHERE is_default = true AND is_active = true;


-- Таблица для платежей
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    job_id BIGINT REFERENCES jobs(id) ON DELETE SET NULL,
    stripe_payment_intent_id VARCHAR(255) NOT NULL UNIQUE,
    stripe_payment_method_id VARCHAR(255) NOT NULL,
    stripe_customer_id VARCHAR(255) NOT NULL,
    amount_cents INT NOT NULL,
    currency VARCHAR(3) DEFAULT 'usd',
    status VARCHAR(50) NOT NULL,
    description TEXT,
    failure_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX idx_user_payment_methods_user_id ON user_payment_methods(user_id);
CREATE INDEX idx_user_payment_methods_stripe_id ON user_payment_methods(stripe_payment_method_id);
CREATE INDEX idx_user_payment_methods_default ON user_payment_methods(user_id, is_default) WHERE is_default = true;

CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_job_id ON payments(job_id);
CREATE INDEX idx_payments_stripe_intent ON payments(stripe_payment_intent_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для автообновления updated_at
CREATE TRIGGER update_user_payment_methods_updated_at 
    BEFORE UPDATE ON user_payment_methods 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_payments_updated_at 
    BEFORE UPDATE ON payments 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();