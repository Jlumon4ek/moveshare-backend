package service

import (
	"context"
	"errors"
	"moveshare/internal/repository"
)

// TruckService defines the interface for truck business logic
type TruckService interface {
	CreateTruck(ctx context.Context, userID int64, truck *repository.Truck) error
	GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error)
	GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error)
	UpdateTruck(ctx context.Context, userID int64, truck *repository.Truck) error
	DeleteTruck(ctx context.Context, userID, truckID int64) error
}

// truckService implements TruckService
type truckService struct {
	truckRepo repository.TruckRepository
}

// NewTruckService creates a new TruckService
func NewTruckService(truckRepo repository.TruckRepository) TruckService {
	return &truckService{truckRepo: truckRepo}
}

// CreateTruck creates a new truck for a user
func (s *truckService) CreateTruck(ctx context.Context, userID int64, truck *repository.Truck) error {
	if truck == nil {
		return errors.New("truck data is required")
	}

	// Validate required fields
	if truck.TruckName == "" {
		return errors.New("truck name is required")
	}
	if truck.LicensePlate == "" {
		return errors.New("license plate is required")
	}
	if truck.Make == "" {
		return errors.New("make is required")
	}
	if truck.Model == "" {
		return errors.New("model is required")
	}
	if truck.Year <= 0 {
		return errors.New("valid year is required")
	}
	if truck.TruckType == "" {
		return errors.New("truck type is required")
	}

	// Validate truck type
	validTypes := map[string]bool{"Small": true, "Medium": true, "Large": true}
	if !validTypes[truck.TruckType] {
		return errors.New("truck type must be Small, Medium, or Large")
	}

	truck.UserID = userID
	return s.truckRepo.CreateTruck(ctx, truck)
}

// GetUserTrucks retrieves all trucks for a user
func (s *truckService) GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error) {
	return s.truckRepo.GetUserTrucks(ctx, userID)
}

// GetTruckByID retrieves a specific truck by ID for a user
func (s *truckService) GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error) {
	return s.truckRepo.GetTruckByID(ctx, userID, truckID)
}

// UpdateTruck updates an existing truck
func (s *truckService) UpdateTruck(ctx context.Context, userID int64, truck *repository.Truck) error {
	if truck == nil {
		return errors.New("truck data is required")
	}

	// Validate truck type if provided
	if truck.TruckType != "" {
		validTypes := map[string]bool{"Small": true, "Medium": true, "Large": true}
		if !validTypes[truck.TruckType] {
			return errors.New("truck type must be Small, Medium, or Large")
		}
	}

	return s.truckRepo.UpdateTruck(ctx, userID, truck)
}

// DeleteTruck deletes a truck
func (s *truckService) DeleteTruck(ctx context.Context, userID, truckID int64) error {
	return s.truckRepo.DeleteTruck(ctx, userID, truckID)
}
