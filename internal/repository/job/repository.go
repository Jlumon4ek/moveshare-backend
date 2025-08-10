package job

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository interface {
	// Job CRUD
	CreateJob(ctx context.Context, job *models.Job, userID int64) error
	GetUserJobs(ctx context.Context, userID int64) ([]models.Job, error)
	GetAvailableJobs(ctx context.Context, userID int64, filters map[string]string, limit, offset int) ([]models.Job, int, error)
	DeleteJob(ctx context.Context, userID, jobID int64) error
	JobExists(ctx context.Context, jobID int64) (bool, error)

	// Job Applications
	ApplyJob(ctx context.Context, userID int64, jobID int64) error
	GetMyApplications(ctx context.Context, userID int64) ([]models.Job, error)
	ChangeJobStatus(ctx context.Context, jobID int64, status string) error

	// Job Photos
	InsertJobPhoto(ctx context.Context, jobID int64, objectID string) error
	GetJobPhotos(ctx context.Context, jobID int64) ([]string, error)
	DeleteJobPhoto(ctx context.Context, jobID int64, photoID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) JobRepository {
	return &repository{db: db}
}
