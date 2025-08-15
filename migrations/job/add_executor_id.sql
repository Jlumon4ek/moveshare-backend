-- Добавляем колонку исполнителя в таблицу jobs
ALTER TABLE jobs ADD COLUMN executor_id BIGINT REFERENCES users(id);

-- Добавляем индекс для быстрого поиска работ по исполнителю
CREATE INDEX idx_jobs_executor_id ON jobs(executor_id);