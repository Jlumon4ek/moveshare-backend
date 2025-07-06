package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TruckType represents the enum values for truck types
type TruckType string

const (
	TruckTypeSmall  TruckType = "Small"
	TruckTypeMedium TruckType = "Medium"
	TruckTypeLarge  TruckType = "Large"
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
	TruckType      TruckType `json:"truck_type"`
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
	// Truck operations
	CreateTruck(ctx context.Context, truck *Truck) error
	GetUserTrucks(ctx context.Context, userID int64) ([]Truck, error)
	GetTruckByID(ctx context.Context, userID, truckID int64) (*Truck, error)
	UpdateTruck(ctx context.Context, userID, truckID int64, truck *Truck) error
	DeleteTruck(ctx context.Context, userID, truckID int64) error
	
	// Photo operations
	AddTruckPhoto(ctx context.Context, photo *TruckPhoto) error
	GetTruckPhotos(ctx context.Context, userID, truckID int64) ([]TruckPhoto, error)
	DeleteTruckPhoto(ctx context.Context, userID, truckID, photoID int64) error
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
		INSERT INTO trucks (user_id, truck_name, license_plate, make, model, year, color,
		                   length, width, height, max_weight, truck_type, climate_control,
		                   liftgate, pallet_jack, security_system, refrigerated, furniture_pads)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query,
		truck.UserID, truck.TruckName, truck.LicensePlate, truck.Make, truck.Model,
		truck.Year, truck.Color, truck.Length, truck.Width, truck.Height,
		truck.MaxWeight, truck.TruckType, truck.ClimateControl, truck.Liftgate,
		truck.PalletJack, truck.SecuritySystem, truck.Refrigerated, truck.FurniturePads,
	).Scan(&truck.ID, &truck.CreatedAt)
	
	return err
}

// GetUserTrucks fetches all trucks belonging to a user
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

	return trucks, rows.Err()
}

// GetTruckByID fetches a specific truck by ID for a user
func (r *truckRepository) GetTruckByID(ctx context.Context, userID, truckID int64) (*Truck, error) {
	query := `
		SELECT id, user_id, truck_name, license_plate, make, model, year, color,
		       length, width, height, max_weight, truck_type, climate_control,
		       liftgate, pallet_jack, security_system, refrigerated, furniture_pads, created_at
		FROM trucks
		WHERE id = $1 AND user_id = $2
	`
	var truck Truck
	err := r.db.QueryRow(ctx, query, truckID, userID).Scan(
		&truck.ID, &truck.UserID, &truck.TruckName, &truck.LicensePlate,
		&truck.Make, &truck.Model, &truck.Year, &truck.Color,
		&truck.Length, &truck.Width, &truck.Height, &truck.MaxWeight,
		&truck.TruckType, &truck.ClimateControl, &truck.Liftgate,
		&truck.PalletJack, &truck.SecuritySystem, &truck.Refrigerated,
		&truck.FurniturePads, &truck.CreatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &truck, nil
}

// UpdateTruck updates an existing truck
func (r *truckRepository) UpdateTruck(ctx context.Context, userID, truckID int64, truck *Truck) error {
	query := `
		UPDATE trucks
		SET truck_name = $3, license_plate = $4, make = $5, model = $6, year = $7, color = $8,
		    length = $9, width = $10, height = $11, max_weight = $12, truck_type = $13,
		    climate_control = $14, liftgate = $15, pallet_jack = $16, security_system = $17,
		    refrigerated = $18, furniture_pads = $19
		WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(ctx, query,
		truckID, userID, truck.TruckName, truck.LicensePlate, truck.Make, truck.Model,
		truck.Year, truck.Color, truck.Length, truck.Width, truck.Height,
		truck.MaxWeight, truck.TruckType, truck.ClimateControl, truck.Liftgate,
		truck.PalletJack, truck.SecuritySystem, truck.Refrigerated, truck.FurniturePads,
	)
	if err != nil {
		return err
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil // No rows updated (truck not found or doesn't belong to user)
	}
	
	return nil
}

// DeleteTruck deletes a truck
func (r *truckRepository) DeleteTruck(ctx context.Context, userID, truckID int64) error {
	query := `DELETE FROM trucks WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, truckID, userID)
	if err != nil {
		return err
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil // No rows deleted (truck not found or doesn't belong to user)
	}
	
	return nil
}

// AddTruckPhoto adds a photo to a truck
func (r *truckRepository) AddTruckPhoto(ctx context.Context, photo *TruckPhoto) error {
	query := `
		INSERT INTO truck_photos (truck_id, user_id, file_name, file_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uploaded_at
	`
	err := r.db.QueryRow(ctx, query,
		photo.TruckID, photo.UserID, photo.FileName, photo.FileURL,
	).Scan(&photo.ID, &photo.UploadedAt)
	
	return err
}

// GetTruckPhotos fetches all photos for a truck
func (r *truckRepository) GetTruckPhotos(ctx context.Context, userID, truckID int64) ([]TruckPhoto, error) {
	// First verify the truck belongs to the user
	truckQuery := `SELECT id FROM trucks WHERE id = $1 AND user_id = $2`
	var truckExists int64
	err := r.db.QueryRow(ctx, truckQuery, truckID, userID).Scan(&truckExists)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil // Truck doesn't exist or doesn't belong to user
		}
		return nil, err
	}

	query := `
		SELECT id, truck_id, user_id, file_name, file_url, uploaded_at
		FROM truck_photos
		WHERE truck_id = $1 AND user_id = $2
		ORDER BY uploaded_at DESC
	`
	rows, err := r.db.Query(ctx, query, truckID, userID)
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

	return photos, rows.Err()
}

// DeleteTruckPhoto deletes a specific truck photo
func (r *truckRepository) DeleteTruckPhoto(ctx context.Context, userID, truckID, photoID int64) error {
	query := `
		DELETE FROM truck_photos 
		WHERE id = $1 AND truck_id = $2 AND user_id = $3
	`
	result, err := r.db.Exec(ctx, query, photoID, truckID, userID)
	if err != nil {
		return err
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil // No rows deleted (photo not found or doesn't belong to user/truck)
	}
	
	return nil
}