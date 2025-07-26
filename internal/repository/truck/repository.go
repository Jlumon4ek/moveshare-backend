package truck

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TruckRepository interface {
	CreateTruck(ctx context.Context, truck *models.Truck) (*models.Truck, error)
	DeleteTruck(ctx context.Context, id int64) error
	GetTruckByID(ctx context.Context, id int64) (*models.Truck, error)
	GetUserTrucks(ctx context.Context, userID int64) ([]*models.Truck, error)
	UpdateTruck(ctx context.Context, userID int64, truck *models.Truck) error
	InsertPhoto(ctx context.Context, truckID int64, objectID string) error
	GetTruckPhotos(ctx context.Context, truckID int64) ([]string, error)
	DeleteTruckPhoto(ctx context.Context, truckID int64, photoID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewTruckRepository(db *pgxpool.Pool) TruckRepository {
	return &repository{db: db}
}
