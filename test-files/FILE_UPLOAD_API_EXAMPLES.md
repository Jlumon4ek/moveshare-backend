# Примеры использования API для загрузки файлов

## Обзор

Теперь загрузка файлов разделена на два типа:
1. **Документы верификации** (`verification_document`) - документы для подтверждения личности/квалификации
2. **Фотографии работы** (`work_photo`) - фотографии выполненной работы

## Новые эндпоинты

### 1. Загрузка документов верификации
```bash
POST /api/jobs/upload-verification-documents/{job_id}
```

**Пример запроса:**
```bash
curl -X POST \
  http://localhost:8080/api/jobs/upload-verification-documents/123 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "files=@verification_doc1.pdf" \
  -F "files=@verification_doc2.jpg"
```

**Ответ:**
```json
{
  "message": "Verification documents uploaded successfully and job status changed to pending",
  "uploaded_files": [
    {
      "id": 0,
      "job_id": 123,
      "file_id": "jobs/123/verification/uuid-verification_doc1.pdf",
      "file_name": "verification_doc1.pdf",
      "file_size": 1048576,
      "content_type": "application/pdf",
      "file_type": "verification_document",
      "uploaded_at": "2024-01-15T10:30:00Z"
    }
  ],
  "files_count": 2,
  "file_type": "verification_document"
}
```

### 2. Загрузка фотографий работы
```bash
POST /api/jobs/upload-work-photos/{job_id}
```

**Пример запроса:**
```bash
curl -X POST \
  http://localhost:8080/api/jobs/upload-work-photos/123 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "files=@work_photo1.jpg" \
  -F "files=@work_photo2.jpg"
```

**Ответ:**
```json
{
  "message": "Work photos uploaded successfully and job status changed to pending",
  "uploaded_files": [
    {
      "id": 0,
      "job_id": 123,
      "file_id": "jobs/123/work-photos/uuid-work_photo1.jpg",
      "file_name": "work_photo1.jpg",
      "file_size": 2097152,
      "content_type": "image/jpeg",
      "file_type": "work_photo",
      "uploaded_at": "2024-01-15T10:35:00Z"
    }
  ],
  "files_count": 2,
  "file_type": "work_photo"
}
```

### 3. Получение файлов по типу
```bash
GET /api/jobs/{job_id}/files/by-type?type={file_type}
```

**Пример запроса для документов верификации:**
```bash
curl -X GET \
  "http://localhost:8080/api/jobs/123/files/by-type?type=verification_document" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Пример запроса для фотографий работы:**
```bash
curl -X GET \
  "http://localhost:8080/api/jobs/123/files/by-type?type=work_photo" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Ответ:**
```json
{
  "job_id": 123,
  "file_type": "verification_document",
  "files": [
    {
      "id": 1,
      "job_id": 123,
      "file_id": "jobs/123/verification/uuid-doc.pdf",
      "file_name": "verification_doc.pdf",
      "file_size": 1048576,
      "content_type": "application/pdf",
      "file_type": "verification_document",
      "uploaded_at": "2024-01-15T10:30:00Z",
      "file_url": "https://minio.example.com/job-files/jobs/123/verification/uuid-doc.pdf?presigned"
    }
  ],
  "count": 1
}
```

### 4. Получение всех файлов (существующий эндпоинт)
```bash
GET /api/jobs/{job_id}/files/
```

**Пример запроса:**
```bash
curl -X GET \
  http://localhost:8080/api/jobs/123/files/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Ответ:**
```json
{
  "job_id": 123,
  "files": [
    {
      "id": 1,
      "job_id": 123,
      "file_id": "jobs/123/verification/uuid-doc.pdf",
      "file_name": "verification_doc.pdf",
      "file_size": 1048576,
      "content_type": "application/pdf",
      "file_type": "verification_document",
      "uploaded_at": "2024-01-15T10:30:00Z",
      "file_url": "https://minio.example.com/presigned-url"
    },
    {
      "id": 2,
      "job_id": 123,
      "file_id": "jobs/123/work-photos/uuid-photo.jpg",
      "file_name": "work_photo.jpg",
      "file_size": 2097152,
      "content_type": "image/jpeg",
      "file_type": "work_photo",
      "uploaded_at": "2024-01-15T10:35:00Z",
      "file_url": "https://minio.example.com/presigned-url"
    }
  ],
  "count": 2
}
```

## Особенности работы

### Разделение файлов по папкам в MinIO:
- **Документы верификации**: `jobs/{job_id}/verification/`
- **Фотографии работы**: `jobs/{job_id}/work-photos/`
- **Старые файлы**: `jobs/{job_id}/` (для обратной совместимости)

### Изменение статуса работы:
- **Документы верификации**: Статус работы изменяется на `pending`
- **Фотографии работы**: Статус работы изменяется на `pending`

### Права доступа:
- **Загрузка файлов**: Только исполнитель работы (executor)
- **Просмотр файлов**: Заказчик (contractor) или исполнитель (executor)

### Валидация статуса:
- Файлы можно загружать только для работ со статусом `claimed` или `pending`

## Типы файлов

| Тип | Описание | Изменяет статус |
|-----|----------|----------------|
| `verification_document` | Документы для верификации | Да (на `pending`) |
| `work_photo` | Фотографии выполненной работы | Да (на `pending`) |
| `legacy` | Старые файлы (для совместимости) | Нет |

## Миграция базы данных

Для поддержки новой функциональности выполните миграцию:
```sql
-- Добавляем колонку file_type
ALTER TABLE job_files ADD COLUMN file_type TEXT;

-- Обновляем существующие записи
UPDATE job_files SET file_type = 'legacy' WHERE file_type IS NULL;

-- Добавляем ограничения
ALTER TABLE job_files ALTER COLUMN file_type SET NOT NULL;
ALTER TABLE job_files ADD CONSTRAINT check_file_type 
CHECK (file_type IN ('verification_document', 'work_photo', 'legacy'));
```