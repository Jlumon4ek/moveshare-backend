package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"moveshare/internal/models"
	"moveshare/internal/repository"
	"moveshare/internal/repository/job"
	"path/filepath"
	"time"
)

type JobService interface {
	// Job CRUD
	CreateJob(ctx context.Context, req *models.CreateJobRequest, userID int64) (*models.CreateJobResponse, error)
	GetUserJobs(ctx context.Context, userID int64) ([]models.JobResponse, error)
	GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.JobResponse, int, error)
	DeleteJob(ctx context.Context, userID, jobID int64) error
	JobExists(ctx context.Context, jobID int64) (bool, error)

	// Job Applications
	ApplyJob(ctx context.Context, userID int64, jobID int64) error
	GetMyApplications(ctx context.Context, userID int64) ([]models.JobResponse, error)

	// Job Photos
	UploadJobPhotos(ctx context.Context, jobID int64, files []*multipart.FileHeader) error
	GetJobPhotos(ctx context.Context, jobID int64) ([]string, error)
}

type jobService struct {
	jobRepo   job.JobRepository
	minioRepo *repository.Repository
}

func NewJobService(jobRepo job.JobRepository, minioRepo *repository.Repository) JobService {
	return &jobService{
		jobRepo:   jobRepo,
		minioRepo: minioRepo,
	}
}

// CreateJob создает новую работу
func (s *jobService) CreateJob(ctx context.Context, req *models.CreateJobRequest, userID int64) (*models.CreateJobResponse, error) {
	// Валидация данных
	if err := s.validateJobRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Парсинг дат
	pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		return nil, fmt.Errorf("invalid pickup date format: %w", err)
	}

	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid delivery date format: %w", err)
	}

	// Расчет расстояния (упрощенная версия)
	distance := s.calculateDistance(req.PickupLocation, req.DeliveryLocation)

	// Расчет общей суммы
	totalAmount := req.CutAmount + req.PaymentAmount

	// Создание модели Job
	job := &models.Job{
		UserID:           userID,
		JobType:          req.JobType,
		JobTitle:         req.JobTitle,
		Description:      req.Description,
		NumberOfBedrooms: req.NumberOfBedrooms,

		PackingBoxes:           req.PackingBoxes,
		BulkyItems:             req.BulkyItems,
		InventoryList:          req.InventoryList,
		Hoisting:               req.Hoisting,
		AdditionalServicesDesc: req.AdditionalServicesDesc,

		TruckSize:      req.TruckSize,
		CrewAssistants: req.CrewAssistants,

		PickupLocation:     req.PickupLocation,
		PickupType:         req.PickupType,
		PickupWalkDistance: req.PickupWalkDistance,

		DeliveryLocation:     req.DeliveryLocation,
		DeliveryType:         req.DeliveryType,
		DeliveryWalkDistance: req.DeliveryWalkDistance,

		PickupDate:        pickupDate,
		PickupTimeStart:   req.PickupTimeStart,
		PickupTimeEnd:     req.PickupTimeEnd,
		DeliveryDate:      deliveryDate,
		DeliveryTimeStart: req.DeliveryTimeStart,
		DeliveryTimeEnd:   req.DeliveryTimeEnd,

		CutAmount:     req.CutAmount,
		PaymentAmount: req.PaymentAmount,
		TotalAmount:   totalAmount,

		Status:        "draft", // Начальный статус
		DistanceMiles: distance,
	}

	// Создание работы в БД
	err = s.jobRepo.CreateJob(ctx, job, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return &models.CreateJobResponse{
		JobID:   job.ID,
		Message: "Job created successfully",
		Status:  job.Status,
		Success: true,
	}, nil
}

// GetUserJobs возвращает все работы пользователя
func (s *jobService) GetUserJobs(ctx context.Context, userID int64) ([]models.JobResponse, error) {
	jobs, err := s.jobRepo.GetUserJobs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user jobs: %w", err)
	}

	var response []models.JobResponse
	for _, job := range jobs {
		// Получаем фотографии для каждой работы
		photoURLs, err := s.getJobPhotoURLs(ctx, job.ID)
		if err != nil {
			// Логируем ошибку, но не останавливаем выполнение
			fmt.Printf("Failed to get photos for job %d: %v\n", job.ID, err)
			photoURLs = []string{}
		}

		jobResponse := s.mapJobToResponse(job)
		jobResponse.PhotoURLs = photoURLs
		response = append(response, jobResponse)
	}

	return response, nil
}

// GetAvailableJobs возвращает доступные работы с фильтрами
func (s *jobService) GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.JobResponse, int, error) {
	jobs, total, err := s.jobRepo.GetAvailableJobs(ctx, userID, filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get available jobs: %w", err)
	}

	var response []models.JobResponse
	for _, job := range jobs {
		// Получаем фотографии для каждой работы
		photoURLs, err := s.getJobPhotoURLs(ctx, job.ID)
		if err != nil {
			// Логируем ошибку, но не останавливаем выполнение
			fmt.Printf("Failed to get photos for job %d: %v\n", job.ID, err)
			photoURLs = []string{}
		}

		jobResponse := s.mapJobToResponse(job)
		jobResponse.PhotoURLs = photoURLs
		response = append(response, jobResponse)
	}

	return response, total, nil
}

// DeleteJob удаляет работу пользователя
func (s *jobService) DeleteJob(ctx context.Context, userID, jobID int64) error {
	return s.jobRepo.DeleteJob(ctx, userID, jobID)
}

// JobExists проверяет существование работы
func (s *jobService) JobExists(ctx context.Context, jobID int64) (bool, error) {
	return s.jobRepo.JobExists(ctx, jobID)
}

// ApplyJob подача заявки на работу
func (s *jobService) ApplyJob(ctx context.Context, userID, jobID int64) error {
	if err := s.jobRepo.ApplyJob(ctx, userID, jobID); err != nil {
		return fmt.Errorf("failed to apply for job: %w", err)
	}

	if err := s.jobRepo.ChangeJobStatus(ctx, jobID, "applied"); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// GetMyApplications возвращает заявки пользователя
func (s *jobService) GetMyApplications(ctx context.Context, userID int64) ([]models.JobResponse, error) {
	jobs, err := s.jobRepo.GetMyApplications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	var response []models.JobResponse
	for _, job := range jobs {
		// Получаем фотографии для каждой работы
		photoURLs, err := s.getJobPhotoURLs(ctx, job.ID)
		if err != nil {
			fmt.Printf("Failed to get photos for job %d: %v\n", job.ID, err)
			photoURLs = []string{}
		}

		jobResponse := s.mapJobToResponse(job)
		jobResponse.PhotoURLs = photoURLs
		response = append(response, jobResponse)
	}

	return response, nil
}

// UploadJobPhotos загружает фотографии для работы
func (s *jobService) UploadJobPhotos(ctx context.Context, jobID int64, files []*multipart.FileHeader) error {
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open photo: %w", err)
		}
		defer file.Close()

		// Читаем данные файла
		data := make([]byte, fileHeader.Size)
		_, err = file.Read(data)
		if err != nil {
			return fmt.Errorf("failed to read photo: %w", err)
		}

		// Генерируем уникальное имя файла
		ext := filepath.Ext(fileHeader.Filename)
		objectName := fmt.Sprintf("job_%d_%d%s", jobID, time.Now().UnixNano(), ext)

		// Загружаем в MinIO
		err = s.minioRepo.UploadBytes(ctx, "job-photos", objectName, data, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("failed to upload photo to MinIO: %w", err)
		}

		// Сохраняем запись в БД
		if err := s.jobRepo.InsertJobPhoto(ctx, jobID, objectName); err != nil {
			return fmt.Errorf("failed to insert photo record: %w", err)
		}
	}

	return nil
}

// GetJobPhotos возвращает URL фотографий работы
func (s *jobService) GetJobPhotos(ctx context.Context, jobID int64) ([]string, error) {
	return s.getJobPhotoURLs(ctx, jobID)
}

// Вспомогательные методы

func (s *jobService) validateJobRequest(req *models.CreateJobRequest) error {
	if req.JobType == "" {
		return fmt.Errorf("job type is required")
	}
	if req.JobTitle == "" {
		return fmt.Errorf("job title is required")
	}
	if req.NumberOfBedrooms == "" {
		return fmt.Errorf("number of bedrooms is required")
	}
	if req.TruckSize == "" {
		return fmt.Errorf("truck size is required")
	}
	if req.CrewAssistants == "" {
		return fmt.Errorf("crew assistants is required")
	}
	if req.PickupLocation == "" {
		return fmt.Errorf("pickup location is required")
	}
	if req.DeliveryLocation == "" {
		return fmt.Errorf("delivery location is required")
	}
	if req.CutAmount < 0 {
		return fmt.Errorf("cut amount cannot be negative")
	}
	if req.PaymentAmount < 0 {
		return fmt.Errorf("payment amount cannot be negative")
	}

	// Валидация дат
	pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		return fmt.Errorf("invalid pickup date format")
	}
	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		return fmt.Errorf("invalid delivery date format")
	}

	// Проверка что даты не в прошлом
	now := time.Now()
	if pickupDate.Before(now.Truncate(24 * time.Hour)) {
		return fmt.Errorf("pickup date cannot be in the past")
	}
	if deliveryDate.Before(pickupDate) {
		return fmt.Errorf("delivery date cannot be before pickup date")
	}

	return nil
}

func (s *jobService) calculateDistance(pickup, delivery string) float64 {
	// Упрощенная версия расчета расстояния
	// В реальном проекте здесь будет использоваться Google Maps API или similar

	// Пока возвращаем фиксированное значение
	return 50.0 // миль
}

func (s *jobService) getJobPhotoURLs(ctx context.Context, jobID int64) ([]string, error) {
	photoIDs, err := s.jobRepo.GetJobPhotos(ctx, jobID)
	if err != nil {
		return nil, err
	}

	var photoURLs []string
	for _, objectName := range photoIDs {
		url, err := s.minioRepo.GetFileURL(ctx, "job-photos", objectName, 10*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("failed to generate URL for %s: %w", objectName, err)
		}
		photoURLs = append(photoURLs, url)
	}

	return photoURLs, nil
}

func (s *jobService) mapJobToResponse(job models.Job) models.JobResponse {
	return models.JobResponse{
		ID:          job.ID,
		UserID:      job.UserID,
		JobType:     job.JobType,
		JobTitle:    job.JobTitle,
		Description: job.Description,

		NumberOfBedrooms: job.NumberOfBedrooms,

		PackingBoxes:           job.PackingBoxes,
		BulkyItems:             job.BulkyItems,
		InventoryList:          job.InventoryList,
		Hoisting:               job.Hoisting,
		AdditionalServicesDesc: job.AdditionalServicesDesc,

		TruckSize:      job.TruckSize,
		CrewAssistants: job.CrewAssistants,

		PickupLocation:     job.PickupLocation,
		PickupType:         job.PickupType,
		PickupWalkDistance: job.PickupWalkDistance,

		DeliveryLocation:     job.DeliveryLocation,
		DeliveryType:         job.DeliveryType,
		DeliveryWalkDistance: job.DeliveryWalkDistance,

		PickupDate:        job.PickupDate,
		PickupTimeStart:   job.PickupTimeStart,
		PickupTimeEnd:     job.PickupTimeEnd,
		DeliveryDate:      job.DeliveryDate,
		DeliveryTimeStart: job.DeliveryTimeStart,
		DeliveryTimeEnd:   job.DeliveryTimeEnd,

		CutAmount:     job.CutAmount,
		PaymentAmount: job.PaymentAmount,
		TotalAmount:   job.TotalAmount,

		Status:        job.Status,
		PhotoURLs:     job.PhotoURLs, // Будет установлено в вызывающем методе
		CreatedAt:     job.CreatedAt,
		UpdatedAt:     job.UpdatedAt,
		DistanceMiles: job.DistanceMiles,
	}
}
