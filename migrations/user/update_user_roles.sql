-- Обновляем существующих пользователей, у которых роль NULL или пустая
UPDATE users 
SET role = 'user' 
WHERE role IS NULL OR role = '';

-- Устанавливаем роль по умолчанию для новых пользователей
ALTER TABLE users 
ALTER COLUMN role SET DEFAULT 'user';

-- Добавляем ограничение для валидных ролей
ALTER TABLE users 
ADD CONSTRAINT check_user_role 
CHECK (role IN ('user', 'admin', 'moderator'));