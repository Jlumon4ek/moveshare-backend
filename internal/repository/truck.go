package repository

import (
	"context"
	"errors"
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
}

// TruckPhoto represents a truck photo entity
type TruckPhoto struct {
	ID         int64     `json:"id"`
	TruckID    int64     `json:"truck_id"`
	UserID     int64     `json:"user_id"`
	FileName   string    `json:"file_name"`
	FileURL    string    `json:"file_url"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// TruckRepository defines the interface for truck data operations
type TruckRepository interface {
	CreateTruck(ctx context.Context, truck *Truck) error
	GetUserTrucks(ctx context.Context, userID int64) ([]Truck, error)
	DeleteTruck(ctx context.Context, userID, truckID int64) error
	GetTruckByID(ctx context.Context, truckID int64) (*Truck, error)
	UploadTruckPhoto(ctx context.Context, photo *TruckPhoto) error
	GetTruckPhotos(ctx context.Context, truckID int64) ([]TruckPhoto, error)
}

// truckRepository implements TruckRepository
type truckRepository struct {
	db *pgxpool.Pool
}

// NewTruckRepository creates a new TruckRepository
func NewTruckRepository(db *pgxpool.Pool) TruckRepository {
	return &truckRepository{db: db}
}

// CreateTruck creates a new truck in the database
func (r *truckRepository) CreateTruck(ctx context.Context, truck *Truck) error {
	query := `
		INSERT INTO trucks (
			user_id, truck_name, license_plate, make, model, year, color,
			length, width, height, max_weight, truck_type, climate_control,
			liftgate, pallet_jack, security_system, refrigerated, furniture_pads
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		truck.UserID, truck.TruckName, truck.LicensePlate, truck.Make, truck.Model,
		truck.Year, truck.Color, truck.Length, truck.Width, truck.Height,
		truck.MaxWeight, truck.TruckType, truck.ClimateControl, truck.Liftgate,
		truck.PalletJack, truck.SecuritySystem, truck.Refrigerated, truck.FurniturePads,
	).Scan(&truck.ID, &truck.CreatedAt)
}

// GetUserTrucks fetches all trucks for a given userID
func (r *truckRepository) GetUserTrucks(ctx context.Context, userID int64) ([]Truck, error) {
	query := `
		SELECT id, user_id, truck_name, license_plate, make, model, year, color,
		       length, width, height, max_weight, truck_type, climate_control,
		       liftgate, pallet_jack, security_system, refrigerated, furniture_pads, created_at
		FROM trucks
		WHERE user_id = $1
		ORDER BY created_at DESC
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
			&truck.ID, &truck.UserID, &truck.TruckName, &truck.LicensePlate,
			&truck.Make, &truck.Model, &truck.Year, &truck.Color,
			&truck.Length, &truck.Width, &truck.Height, &truck.MaxWeight,
			&truck.TruckType, &truck.ClimateControl, &truck.Liftgate,
			&truck.PalletJack, &truck.SecuritySystem, &truck.Refrigerated,
			&truck.FurniturePads, &truck.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		trucks = append(trucks, truck)
	}

	return trucks, nil
}

// DeleteTruck deletes a truck owned by the user
func (r *truckRepository) DeleteTruck(ctx context.Context, userID, truckID int64) error {
	query := `DELETE FROM trucks WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, truckID, userID)
	if err != nil {
		return err
	}
	
	if cmdTag.RowsAffected() == 0 {
		return ErrTruckNotFound
	}
	
	return nil
}

// GetTruckByID fetches a truck by its ID
func (r *truckRepository) GetTruckByID(ctx context.Context, truckID int64) (*Truck, error) {
	query := `
		SELECT id, user_id, truck_name, license_plate, make, model, year, color,
		       length, width, height, max_weight, truck_type, climate_control,
		       liftgate, pallet_jack, security_system, refrigerated, furniture_pads, created_at
		FROM trucks
		WHERE id = $1
	`
	var truck Truck
	err := r.db.QueryRow(ctx, query, truckID).Scan(
		&truck.ID, &truck.UserID, &truck.TruckName, &truck.LicensePlate,
		&truck.Make, &truck.Model, &truck.Year, &truck.Color,
		&truck.Length, &truck.Width, &truck.Height, &truck.MaxWeight,
		&truck.TruckType, &truck.ClimateControl, &truck.Liftgate,
		&truck.PalletJack, &truck.SecuritySystem, &truck.Refrigerated,
		&truck.FurniturePads, &truck.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &truck, nil
}

// UploadTruckPhoto uploads a photo for a truck
func (r *truckRepository) UploadTruckPhoto(ctx context.Context, photo *TruckPhoto) error {
	query := `
		INSERT INTO truck_photos (truck_id, user_id, file_name, file_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uploaded_at
	`
	return r.db.QueryRow(ctx, query,
		photo.TruckID, photo.UserID, photo.FileName, photo.FileURL,
	).Scan(&photo.ID, &photo.UploadedAt)
}

// GetTruckPhotos fetches all photos for a truck
func (r *truckRepository) GetTruckPhotos(ctx context.Context, truckID int64) ([]TruckPhoto, error) {
	query := `
		SELECT id, truck_id, user_id, file_name, file_url, uploaded_at
		FROM truck_photos
		WHERE truck_id = $1
		ORDER BY uploaded_at DESC
	`
	rows, err := r.db.Query(ctx, query, truckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []TruckPhoto
	for rows.Next() {
		var photo TruckPhoto
		err := rows.Scan(
			&photo.ID, &photo.TruckID, &photo.UserID,
			&photo.FileName, &photo.FileURL, &photo.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

// Custom errors
var (
	ErrTruckNotFound = errors.New("truck not found")
)