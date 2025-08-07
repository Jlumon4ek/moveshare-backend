package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository/job"
)

type JobService interface {
	ApplyJob(ctx context.Context, userID int64, jobID int64) error
	CreateJob(ctx context.Context, job *models.Job, userId int64) error
	DeleteJob(ctx context.Context, userID, jobID int64) error
	GetMyApplications(ctx context.Context, userID int64) ([]models.Job, error)
	GetUserJobs(ctx context.Context, userID int64) ([]models.Job, error)

	GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.Job, int64, error)
}

type jobService struct {
	jobRepo job.JobRepository
}

func NewJobService(jobRepo job.JobRepository) JobService {
	return &jobService{
		jobRepo: jobRepo,
	}
}

func (s *jobService) ApplyJob(ctx context.Context, userID, jobID int64) error {
	if err := s.jobRepo.ApplyJob(ctx, userID, jobID); err != nil {
		return err
	}

	if err := s.jobRepo.ChangeJobStatus(ctx, jobID, "applied"); err != nil {
		return err
	}

	return nil
}

func (s *jobService) CreateJob(ctx context.Context, job *models.Job, userId int64) error {
	return s.jobRepo.CreateJob(ctx, job, userId)
}

func (s *jobService) DeleteJob(ctx context.Context, userID int64, jobID int64) error {
	return s.jobRepo.DeleteJob(ctx, userID, jobID)
}

func (s *jobService) GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.Job, int64, error) {
	total, err := s.jobRepo.GetTotalJobCount(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	jobs, err := s.jobRepo.GetAvailableJobs(ctx, userID, filters, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (s *jobService) GetMyApplications(ctx context.Context, userID int64) ([]models.Job, error) {
	return s.jobRepo.GetMyApplications(ctx, userID)
}

func (s *jobService) GetUserJobs(ctx context.Context, userID int64) ([]models.Job, error) {
	return s.jobRepo.GetUserJobs(ctx, userID)
}
