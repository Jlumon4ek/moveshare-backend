package service

import (
	"context"
	"errors"
	"moveshare/internal/repository"
)

// JobService defines the interface for job business logic
type JobService interface {
	CreateJob(ctx context.Context, userID int64, job *repository.Job) error
	GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]repository.Job, error)
	GetUserJobs(ctx context.Context, userID int64) ([]repository.Job, error)
	DeleteJob(ctx context.Context, userID, jobID int64) error
	ApplyForJob(ctx context.Context, userID, jobID int64) error
	GetMyApplications(ctx context.Context, userID int64) ([]repository.Job, error) // New method
}

// jobService implements JobService
type jobService struct {
	jobRepo repository.JobRepository
}

// NewJobService creates a new JobService
func NewJobService(jobRepo repository.JobRepository) JobService {
	return &jobService{jobRepo: jobRepo}
}

// CreateJob creates a new job for a user
func (s *jobService) CreateJob(ctx context.Context, userID int64, job *repository.Job) error {
	if job.JobTitle == "" || job.PickupLocation == "" || job.DeliveryLocation == "" {
		return errors.New("job title, pickup location, and delivery location are required")
	}
	job.UserID = userID
	return s.jobRepo.CreateJob(ctx, job)
}

// GetAvailableJobs fetches jobs excluding those created by the given userID with filters and pagination
func (s *jobService) GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]repository.Job, error) {
	return s.jobRepo.GetAvailableJobs(ctx, userID, filters, limit, offset)
}

// GetUserJobs fetches all jobs created by the given userID
func (s *jobService) GetUserJobs(ctx context.Context, userID int64) ([]repository.Job, error) {
	return s.jobRepo.GetUserJobs(ctx, userID)
}

// DeleteJob deletes a job if it belongs to the given userID
func (s *jobService) DeleteJob(ctx context.Context, userID, jobID int64) error {
	return s.jobRepo.DeleteJob(ctx, userID, jobID)
}

// ApplyForJob allows a user to apply for a job
func (s *jobService) ApplyForJob(ctx context.Context, userID, jobID int64) error {
	// Check if the user is not the job creator
	jobs, err := s.jobRepo.GetUserJobs(ctx, userID)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if job.ID == jobID {
			return errors.New("cannot apply to your own job")
		}
	}

	return s.jobRepo.ApplyForJob(ctx, userID, jobID)
}

// GetMyApplications fetches all applications submitted by the given userID
func (s *jobService) GetMyApplications(ctx context.Context, userID int64) ([]repository.Job, error) {
	// This is a placeholder; implement the actual logic in the repository
	return s.jobRepo.GetMyApplications(ctx, userID)
}
