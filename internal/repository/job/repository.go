package job

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository interface {
	ApplyJob(ctx context.Context, userID, jobID int64) error
	CreateJob(ctx context.Context, job *models.Job) error
	DeleteJob(ctx context.Context, userID, jobID int64) error
	GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.Job, error)
	GetMyApplications(ctx context.Context, userID int64) ([]models.Job, error)
	GetUserJobs(ctx context.Context, userID int64) ([]models.Job, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) JobRepository {
	return &repository{db: db}
}
