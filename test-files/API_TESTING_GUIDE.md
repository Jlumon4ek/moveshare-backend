# Job Files Upload API Testing Guide

## Описание функциональности
Система позволяет загружать файлы к claimed job'ам и автоматически меняет статус на "pending".

## Структура таблиц
- `jobs` - основная таблица работ
- `job_files` - таблица файлов (по аналогии с truck_photos)

## API Endpoints

### 1. Логин (получение JWT токена)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-password"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIs...",
  "user": { ... }
}
```

### 2. Создание job (если нужно для теста)
```bash
curl -X POST http://localhost:8080/api/jobs/post-new-job/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "Moving",
    "truck_size": "Large",
    "pickup_address": "123 Main St, City",
    "delivery_address": "456 Oak Ave, City",
    "pickup_date": "2025-01-20",
    "pickup_time_from": "09:00",
    "pickup_time_to": "10:00",
    "delivery_date": "2025-01-20",
    "delivery_time_from": "14:00",
    "delivery_time_to": "15:00",
    "payment_amount": 500.00
  }'
```

### 3. Claim job
```bash
curl -X POST http://localhost:8080/api/jobs/claim-job/JOB_ID/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "message": "Job claimed successfully"
}
```

### 4. 🚀 ОСНОВНАЯ ФУНКЦИОНАЛЬНОСТЬ: Upload files
```bash
curl -X POST http://localhost:8080/api/jobs/upload-files/JOB_ID/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "files=@test-files/test-document.txt" \
  -F "files=@test-files/inventory.csv"
```

**Response:**
```json
{
  "message": "Files uploaded successfully and job status changed to pending",
  "uploaded_files": [
    {
      "id": 0,
      "job_id": 123,
      "file_id": "jobs/123/uuid-filename.txt",
      "file_name": "test-document.txt",
      "file_size": 1024,
      "content_type": "text/plain",
      "uploaded_at": "2025-01-16T10:30:00Z"
    }
  ],
  "files_count": 2
}
```

### 5. Получить claimed jobs с файлами
```bash
curl -X GET http://localhost:8080/api/jobs/claimed-jobs/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "jobs": [
    {
      "id": 123,
      "job_status": "pending",
      "files": [
        {
          "id": 1,
          "job_id": 123,
          "file_id": "jobs/123/uuid-filename.txt",
          "file_name": "test-document.txt",
          "file_url": "https://minio-url/presigned-url",
          "uploaded_at": "2025-01-16T10:30:00Z"
        }
      ]
    }
  ]
}
```

## Postman Collection Setup

### Environment Variables
```
jwt_token = your_jwt_token_here
base_url = http://localhost:8080/api
job_id = your_job_id_here
```

### 1. Login Request
- **Method**: POST
- **URL**: `{{base_url}}/auth/login`
- **Body** (JSON):
```json
{
  "email": "your-email@example.com",
  "password": "your-password"
}
```
- **Tests** (для автоматического сохранения токена):
```javascript
if (pm.response.code === 200) {
    const response = pm.response.json();
    pm.environment.set("jwt_token", response.token);
}
```

### 2. Upload Files Request
- **Method**: POST
- **URL**: `{{base_url}}/jobs/upload-files/{{job_id}}/`
- **Headers**: 
  - `Authorization: Bearer {{jwt_token}}`
- **Body**: 
  - Type: `form-data`
  - Key: `files` (type: File) - выберите файлы
  - Key: `files` (type: File) - можно добавить несколько файлов

### 3. Get Claimed Jobs
- **Method**: GET
- **URL**: `{{base_url}}/jobs/claimed-jobs/`
- **Headers**: 
  - `Authorization: Bearer {{jwt_token}}`

## Бизнес-логика

### Workflow:
1. Пользователь создает job → статус `open`
2. Другой пользователь делает claim → статус `claimed`
3. Исполнитель загружает файлы → статус автоматически меняется на `pending`
4. Файлы доступны через `/claimed-jobs/` с presigned URLs

### Ограничения:
- Только исполнитель (executor) может загружать файлы
- Job должен быть в статусе `claimed`
- После загрузки статус меняется на `pending`

### Хранение файлов:
- **MinIO bucket**: `job-files`
- **Path pattern**: `jobs/{job_id}/{uuid}-{filename}`
- **Presigned URLs**: действительны 24 часа

## Тестовые сценарии

### Позитивные тесты:
1. ✅ Загрузка одного файла
2. ✅ Загрузка нескольких файлов
3. ✅ Проверка смены статуса на pending
4. ✅ Получение файлов через /claimed-jobs/

### Негативные тесты:
1. ❌ Загрузка без токена → 401
2. ❌ Загрузка к чужому job → 403
3. ❌ Загрузка к job в неправильном статусе → 400
4. ❌ Загрузка без файлов → 400

## Проверка результатов

### В БД:
```sql
-- Проверить файлы
SELECT * FROM job_files WHERE job_id = YOUR_JOB_ID;

-- Проверить статус job
SELECT id, job_status FROM jobs WHERE id = YOUR_JOB_ID;
```

### В MinIO:
- Bucket: `job-files`
- Objects: `jobs/{job_id}/...`

---

**Примечание**: Замените `YOUR_JWT_TOKEN` и `JOB_ID` на реальные значения из ваших тестов.