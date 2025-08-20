package notifications

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.NotificationRequest) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int, typeFilter string, unreadOnly bool) ([]models.Notification, int, error)
	GetByID(ctx context.Context, id, userID int64) (*models.Notification, error)
	MarkAsRead(ctx context.Context, id, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	Delete(ctx context.Context, id, userID int64) error
	DeleteAll(ctx context.Context, userID int64) error
	GetStats(ctx context.Context, userID int64) (*models.NotificationStats, error)
	CleanupExpired(ctx context.Context) (int, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository {
	return &repository{db: db}
}