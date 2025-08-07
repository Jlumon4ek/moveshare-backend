package job

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository interface {
	ApplyJob(ctx context.Context, userID int64, jobID int64) error
	CreateJob(ctx context.Context, job *models.Job, userId int64) error
	DeleteJob(ctx context.Context, userID int64, jobID int64) error
	GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit int, offset int) ([]models.Job, error)
	GetMyApplications(ctx context.Context, userID int64) ([]models.Job, error)
	GetUserJobs(ctx context.Context, userID int64) ([]models.Job, error)
	ChangeJobStatus(ctx context.Context, jobID int64, status string) error
	GetTotalJobCount(ctx context.Context, userID int64) (int64, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) JobRepository {
	return &repository{db: db}
}
