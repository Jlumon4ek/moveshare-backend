package service

import (
	"context"
	"fmt"
	"moveshare/internal/config"
	"moveshare/internal/models"
	"moveshare/internal/repository"
	"moveshare/internal/utils"
	"time"
)

type JobService struct {
	jobRepo       *repository.JobRepository
	googleMapsCfg *config.GoogleMapsConfig
	minioRepo     *repository.Repository
}

func NewJobService(jobRepo *repository.JobRepository, googleMapsCfg *config.GoogleMapsConfig, minioRepo *repository.Repository) *JobService {
	return &JobService{
		jobRepo:       jobRepo,
		googleMapsCfg: googleMapsCfg,
		minioRepo:     minioRepo,
	}
}

func (s *JobService) CreateJob(userID int64, req *models.CreateJobRequest) (*models.Job, error) {
	pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		return nil, err
	}

	pickupTimeFrom, err := time.Parse("15:04", req.PickupTimeFrom)
	if err != nil {
		return nil, err
	}

	pickupTimeTo, err := time.Parse("15:04", req.PickupTimeTo)
	if err != nil {
		return nil, err
	}

	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		return nil, err
	}

	deliveryTimeFrom, err := time.Parse("15:04", req.DeliveryTimeFrom)
	if err != nil {
		return nil, err
	}

	deliveryTimeTo, err := time.Parse("15:04", req.DeliveryTimeTo)
	if err != nil {
		return nil, err
	}

	// Calculate distance using Google Maps API
	fmt.Printf("Calculating distance from '%s' to '%s'\n", req.PickupAddress, req.DeliveryAddress)
	distanceResult, err := utils.GetDistanceFromAddresses(req.PickupAddress, req.DeliveryAddress, s.googleMapsCfg)
	if err != nil {
		fmt.Printf("ERROR: Failed to calculate distance: %v\n", err)
		// Use provided distance or default to 0
		if req.DistanceMiles == 0 {
			req.DistanceMiles = 0
		}
	} else {
		// Convert meters to miles (1 meter = 0.000621371 miles)
		req.DistanceMiles = float64(distanceResult.DistanceValue) * 0.000621371
		fmt.Printf("SUCCESS: Distance calculated: %s (%d meters) = %.2f miles\n", 
			distanceResult.Distance, distanceResult.DistanceValue, req.DistanceMiles)
	}

	job := &models.Job{
		ContractorID:                  userID,
		JobType:                       req.JobType,
		NumberOfBedrooms:              req.NumberOfBedrooms,
		PackingBoxes:                  req.PackingBoxes,
		BulkyItems:                    req.BulkyItems,
		InventoryList:                 req.InventoryList,
		Hoisting:                      req.Hoisting,
		AdditionalServicesDescription: req.AdditionalServicesDescription,
		EstimatedCrewAssistants:       req.EstimatedCrewAssistants,
		TruckSize:                     req.TruckSize,
		PickupAddress:                 req.PickupAddress,
		PickupFloor:                   req.PickupFloor,
		PickupBuildingType:            req.PickupBuildingType,
		PickupWalkDistance:            req.PickupWalkDistance,
		DeliveryAddress:               req.DeliveryAddress,
		DeliveryFloor:                 req.DeliveryFloor,
		DeliveryBuildingType:          req.DeliveryBuildingType,
		DeliveryWalkDistance:          req.DeliveryWalkDistance,
		DistanceMiles:                 req.DistanceMiles,
		JobStatus:                     "active",
		PickupDate:                    pickupDate,
		PickupTimeFrom:                pickupTimeFrom,
		PickupTimeTo:                  pickupTimeTo,
		DeliveryDate:                  deliveryDate,
		DeliveryTimeFrom:              deliveryTimeFrom,
		DeliveryTimeTo:                deliveryTimeTo,
		CutAmount:                     req.CutAmount,
		PaymentAmount:                 req.PaymentAmount,
		WeightLbs:                     req.WeightLbs,
		VolumeCuFt:                    req.VolumeCuFt,
	}

	ctx := context.Background()
	err = s.jobRepo.CreateJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// internal/service/job.go - обновить метод GetAvailableJobs

func (s *JobService) GetAvailableJobs(userID int64, filters *models.JobFilters) ([]models.AvailableJobDTO, int, error) {
	ctx := context.Background()

	// Валидация фильтров
	if err := filters.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid filters: %w", err)
	}

	// Получаем задания с фильтрацией
	jobs, total, err := s.jobRepo.GetAvailableJobs(ctx, userID, filters)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (s *JobService) GetFilterOptions(userID int64) (*models.JobFilterOptions, error) {
	ctx := context.Background()

	return s.jobRepo.GetFilterOptions(ctx, userID)
}

func (s *JobService) GetJobByID(jobID int64) (*models.Job, error) {
	ctx := context.Background()
	return s.jobRepo.GetJobByID(ctx, jobID)
}

func (s *JobService) DeleteJob(jobID, userID int64) error {
	ctx := context.Background()
	return s.jobRepo.DeleteJob(ctx, jobID, userID)
}

func (s *JobService) ClaimJob(jobID, userID int64) error {
	ctx := context.Background()
	return s.jobRepo.ClaimJob(ctx, jobID, userID)
}

func (s *JobService) GetMyJobs(userID int64, page, limit int) ([]models.Job, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	jobs, err := s.jobRepo.GetMyJobs(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.jobRepo.GetCountMyJobs(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (s *JobService) JobExists(jobID int64) (bool, error) {
	ctx := context.Background()
	return s.jobRepo.JobExists(ctx, jobID)
}

func (s *JobService) GetClaimedJobs(userID int64, page, limit int) ([]models.Job, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	jobs, err := s.jobRepo.GetClaimedJobs(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Получаем файлы для каждой работы
	for i := range jobs {
		files, err := s.GetJobFiles(jobs[i].ID)
		if err != nil {
			// Логируем ошибку, но не прерываем работу
			fmt.Printf("Failed to get files for job %d: %v\n", jobs[i].ID, err)
			continue
		}
		jobs[i].Files = files
	}

	total, err := s.jobRepo.GetCountClaimedJobs(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (s *JobService) MarkJobCompleted(jobID, userID int64) error {
	ctx := context.Background()
	return s.jobRepo.MarkJobCompleted(ctx, jobID, userID)
}

func (s *JobService) GetJobsForExport(userID int64, jobIDs []int64) ([]models.Job, error) {
	ctx := context.Background()
	return s.jobRepo.GetJobsByIDs(ctx, userID, jobIDs)
}

func (s *JobService) GetJobsStats(userID int64) (models.JobsStats, error) {
	ctx := context.Background()
	return s.jobRepo.GetJobsStats(ctx, userID)
}


func (s *JobService) GetUserWorkStats(userID int64) (models.UserWorkStats, error) {
	ctx := context.Background()
	return s.jobRepo.GetUserWorkStats(ctx, userID)
}

func (s *JobService) GetPendingJobs(userID int64, limit int) ([]models.Job, error) {
	ctx := context.Background()
	jobs, err := s.jobRepo.GetPendingJobs(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *JobService) GetTodayScheduleJobs(userID int64, page, limit int) ([]models.Job, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	jobs, err := s.jobRepo.GetTodayScheduleJobs(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.jobRepo.GetCountTodayScheduleJobs(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (s *JobService) UploadJobFile(jobID int64, fileID, fileName string, fileSize int64, contentType string) error {
	ctx := context.Background()
	return s.jobRepo.InsertJobFile(ctx, jobID, fileID, fileName, fileSize, contentType)
}

func (s *JobService) GetJobFiles(jobID int64) ([]models.JobFile, error) {
	ctx := context.Background()
	files, err := s.jobRepo.GetJobFiles(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Генерируем URL для каждого файла
	for i := range files {
		fileURL, err := s.minioRepo.GetFileURL(ctx, "job-files", files[i].FileID, 24*time.Hour)
		if err != nil {
			// Если не можем получить URL, логируем ошибку но не прерываем
			fmt.Printf("Failed to get URL for file %s: %v\n", files[i].FileID, err)
			continue
		}
		files[i].FileURL = fileURL
	}

	return files, nil
}

func (s *JobService) MarkJobAsPending(jobID int64) error {
	ctx := context.Background()
	return s.jobRepo.UpdateJobStatus(ctx, jobID, "pending")
}

func (s *JobService) UploadJobFileWithType(jobID int64, fileID, fileName string, fileSize int64, contentType, fileType string) error {
	ctx := context.Background()
	return s.jobRepo.InsertJobFileWithType(ctx, jobID, fileID, fileName, fileSize, contentType, fileType)
}

func (s *JobService) GetJobFilesByType(jobID int64, fileType string) ([]models.JobFile, error) {
	ctx := context.Background()
	files, err := s.jobRepo.GetJobFilesByType(ctx, jobID, fileType)
	if err != nil {
		return nil, err
	}

	// Генерируем URL для каждого файла
	for i := range files {
		fileURL, err := s.minioRepo.GetFileURL(ctx, "job-files", files[i].FileID, 24*time.Hour)
		if err != nil {
			// Если не можем получить URL, логируем ошибку но не прерываем
			fmt.Printf("Failed to get URL for file %s: %v\n", files[i].FileID, err)
			continue
		}
		files[i].FileURL = fileURL
	}

	return files, nil
}
