package admin

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository interface {
	GetUserCount(ctx context.Context) (int, error)
	GetChatConversationCount(ctx context.Context) (int, error)
	GetUsersList(ctx context.Context, limit, offset int) ([]models.UserCompanyInfo, error)
	GetAllJobs(ctx context.Context, limit, offset int) ([]models.Job, error)
	ChangeUserStatus(ctx context.Context, userID int, newStatus string) error
	GetUserRole(ctx context.Context, userID int64) (string, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) AdminRepository {
	return &repository{db: db}
}
