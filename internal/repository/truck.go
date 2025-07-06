package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Truck represents a truck entity
type Truck struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	TruckName      string    `json:"truck_name"`
	LicensePlate   string    `json:"license_plate"`
	Make           string    `json:"make"`
	Model          string    `json:"model"`
	Year           int       `json:"year"`
	Color          string    `json:"color"`
	Length         float64   `json:"length"`
	Width          float64   `json:"width"`
	Height         float64   `json:"height"`
	MaxWeight      float64   `json:"max_weight"`
	TruckType      string    `json:"truck_type"`
	ClimateControl bool      `json:"climate_control"`
	Liftgate       bool      `json:"liftgate"`
	PalletJack     bool      `json:"pallet_jack"`
	SecuritySystem bool      `json:"security_system"`
	Refrigerated   bool      `json:"refrigerated"`
	FurniturePads  bool      `json:"furniture_pads"`
	CreatedAt      time.Time `json:"created_at"`
	Photos         []string  `json:"photos"` // URLs of uploaded photos
}

// TruckRepository defines the interface for truck data operations
type TruckRepository interface {
	GetUserTrucks(ctx context.Context, userID int64) ([]Truck, error)
	CreateTruck(ctx context.Context, truck *Truck, photoURLs []string) (*Truck, error)
}

// truckRepository implements TruckRepository
type truckRepository struct {
	db *pgxpool.Pool
}

// NewTruckRepository creates a new TruckRepository
func NewTruckRepository(db *pgxpool.Pool) TruckRepository {
	return &truckRepository{db: db}
}

// GetUserTrucks fetches all trucks for the given userID
func (r *truckRepository) GetUserTrucks(ctx context.Context, userID int64) ([]Truck, error) {
	query := `
		SELECT t.id, t.user_id, t.truck_name, t.license_plate, t.make, t.model, t.year, t.color,
		       t.length, t.width, t.height, t.max_weight, t.truck_type, t.climate_control,
		       t.liftgate, t.pallet_jack, t.security_system, t.refrigerated, t.furniture_pads,
		       t.created_at
		FROM trucks t
		WHERE t.user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trucks []Truck
	for rows.Next() {
		var truck Truck
		err := rows.Scan(
			&truck.ID, &truck.UserID, &truck.TruckName, &truck.LicensePlate, &truck.Make,
			&truck.Model, &truck.Year, &truck.Color, &truck.Length, &truck.Width, &truck.Height,
			&truck.MaxWeight, &truck.TruckType, &truck.ClimateControl, &truck.Liftgate,
			&truck.PalletJack, &truck.SecuritySystem, &truck.Refrigerated, &truck.FurniturePads,
			&truck.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Fetch photos for this truck
		photos, err := r.getTruckPhotos(ctx, truck.ID, userID)
		if err != nil {
			return nil, err
		}
		truck.Photos = photos
		trucks = append(trucks, truck)
	}

	return trucks, nil
}

// CreateTruck creates a new truck and associates photos
func (r *truckRepository) CreateTruck(ctx context.Context, truck *Truck, photoURLs []string) (*Truck, error) {
	query := `
		INSERT INTO trucks (user_id, truck_name, license_plate, make, model, year, color, length,
		                   width, height, max_weight, truck_type, climate_control, liftgate,
		                   pallet_jack, security_system, refrigerated, furniture_pads)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query,
		truck.UserID, truck.TruckName, truck.LicensePlate, truck.Make, truck.Model, truck.Year,
		truck.Color, truck.Length, truck.Width, truck.Height, truck.MaxWeight, truck.TruckType,
		truck.ClimateControl, truck.Liftgate, truck.PalletJack, truck.SecuritySystem,
		truck.Refrigerated, truck.FurniturePads,
	).Scan(&truck.ID, &truck.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Insert photos
	for _, url := range photoURLs {
		query := `INSERT INTO truck_photos (truck_id, user_id, file_name, file_url) VALUES ($1, $2, $3, $3)`
		_, err := r.db.Exec(ctx, query, truck.ID, truck.UserID, url)
		if err != nil {
			return nil, err
		}
	}

	// Fetch photos to include in the returned truck
	photos, err := r.getTruckPhotos(ctx, truck.ID, truck.UserID)
	if err != nil {
		return nil, err
	}
	truck.Photos = photos

	return truck, nil
}

// getTruckPhotos fetches photo URLs for a truck
func (r *truckRepository) getTruckPhotos(ctx context.Context, truckID, userID int64) ([]string, error) {
	query := `SELECT file_url FROM truck_photos WHERE truck_id = $1 AND user_id = $2`
	rows, err := r.db.Query(ctx, query, truckID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		photos = append(photos, url)
	}
	return photos, nil
}
