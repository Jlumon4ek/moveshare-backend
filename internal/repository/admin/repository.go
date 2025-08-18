package admin

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository interface {
	GetUserCount(ctx context.Context) (int, error)
	GetPendingUsersCount(ctx context.Context) (int, error)
	GetChatConversationCount(ctx context.Context) (int, error)
	GetActiveJobsCount(ctx context.Context) (int, error)
	GetUsersList(ctx context.Context, limit, offset int) ([]models.UserCompanyInfo, error)
	GetUsersListTotal(ctx context.Context) (int, error)
	GetJobsList(ctx context.Context, limit, offset int, statuses []string) ([]models.JobManagementInfo, error)
	GetJobsListTotal(ctx context.Context, statuses []string) (int, error)
	// GetAllJobs(ctx context.Context, limit, offset int) ([]models.Job, error)
	ChangeUserStatus(ctx context.Context, userID int, newStatus string) error
	GetUserRole(ctx context.Context, userID int64) (string, error)
	ChangeVerificationFileStatus(ctx context.Context, fileID int, newStatus string) error
	GetUserFullInfo(ctx context.Context, userID int64) (*models.UserFullInfo, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) AdminRepository {
	return &repository{db: db}
}
