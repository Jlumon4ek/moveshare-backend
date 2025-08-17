-- Добавляем колонку file_type в существующую таблицу job_files
ALTER TABLE job_files 
ADD COLUMN file_type TEXT;

-- Обновляем существующие записи, присваивая им тип 'legacy'
UPDATE job_files 
SET file_type = 'legacy' 
WHERE file_type IS NULL;

-- Добавляем ограничение после обновления данных
ALTER TABLE job_files 
ALTER COLUMN file_type SET NOT NULL;

-- Добавляем проверочное ограничение
ALTER TABLE job_files 
ADD CONSTRAINT check_file_type 
CHECK (file_type IN ('verification_document', 'work_photo', 'legacy'));