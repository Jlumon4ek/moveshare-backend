package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"moveshare/internal/repository"
)

// TruckService defines the interface for truck business logic
type TruckService interface {
	CreateTruck(ctx context.Context, userID int64, truck *repository.Truck) error
	GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error)
	GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error)
	UpdateTruck(ctx context.Context, userID int64, truck *repository.Truck) error
	DeleteTruck(ctx context.Context, userID, truckID int64) error
	
	UploadTruckPhotos(ctx context.Context, userID, truckID int64, files []*multipart.FileHeader) error
	GetTruckPhotos(ctx context.Context, userID, truckID int64) ([]repository.TruckPhoto, error)
}

// truckService implements TruckService
type truckService struct {
	truckRepo repository.TruckRepository
}

// NewTruckService creates a new TruckService
func NewTruckService(truckRepo repository.TruckRepository) TruckService {
	return &truckService{truckRepo: truckRepo}
}

// CreateTruck creates a new truck with validation
func (s *truckService) CreateTruck(ctx context.Context, userID int64, truck *repository.Truck) error {
	if err := s.validateTruck(truck); err != nil {
		return err
	}
	
	truck.UserID = userID
	return s.truckRepo.CreateTruck(ctx, truck)
}

// GetUserTrucks fetches all trucks for a user
func (s *truckService) GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error) {
	return s.truckRepo.GetUserTrucks(ctx, userID)
}

// GetTruckByID fetches a truck by ID with ownership validation
func (s *truckService) GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error) {
	return s.truckRepo.GetTruckByID(ctx, userID, truckID)
}

// UpdateTruck updates a truck with validation
func (s *truckService) UpdateTruck(ctx context.Context, userID int64, truck *repository.Truck) error {
	if err := s.validateTruck(truck); err != nil {
		return err
	}
	
	truck.UserID = userID
	return s.truckRepo.UpdateTruck(ctx, truck)
}

// DeleteTruck deletes a truck with ownership validation
func (s *truckService) DeleteTruck(ctx context.Context, userID, truckID int64) error {
	// First verify ownership
	_, err := s.truckRepo.GetTruckByID(ctx, userID, truckID)
	if err != nil {
		return errors.New("truck not found or access denied")
	}
	
	return s.truckRepo.DeleteTruck(ctx, userID, truckID)
}

// UploadTruckPhotos handles truck photo uploads with validation
func (s *truckService) UploadTruckPhotos(ctx context.Context, userID, truckID int64, files []*multipart.FileHeader) error {
	// Verify truck ownership
	_, err := s.truckRepo.GetTruckByID(ctx, userID, truckID)
	if err != nil {
		return errors.New("truck not found or access denied")
	}
	
	// Check current photo count
	currentCount, err := s.truckRepo.GetTruckPhotoCount(ctx, truckID)
	if err != nil {
		return fmt.Errorf("failed to check current photo count: %w", err)
	}
	
	// Validate photo count (max 10 total)
	if currentCount+len(files) > 10 {
		return fmt.Errorf("cannot upload %d photos, maximum 10 photos per truck (currently have %d)", len(files), currentCount)
	}
	
	// Validate each file
	for _, file := range files {
		if err := s.validatePhotoFile(file); err != nil {
			return err
		}
	}
	
	// Upload each file
	for _, file := range files {
		photo := &repository.TruckPhoto{
			TruckID:  truckID,
			UserID:   userID,
			FileName: file.Filename,
			FileURL:  fmt.Sprintf("/uploads/trucks/%d/%s", truckID, file.Filename), // Simple file URL
		}
		
		if err := s.truckRepo.AddTruckPhoto(ctx, photo); err != nil {
			return fmt.Errorf("failed to save photo %s: %w", file.Filename, err)
		}
	}
	
	return nil
}

// GetTruckPhotos fetches all photos for a truck
func (s *truckService) GetTruckPhotos(ctx context.Context, userID, truckID int64) ([]repository.TruckPhoto, error) {
	return s.truckRepo.GetTruckPhotos(ctx, userID, truckID)
}

// validateTruck validates truck data
func (s *truckService) validateTruck(truck *repository.Truck) error {
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
		return errors.New("year is required and must be positive")
	}
	
	if truck.Color == "" {
		return errors.New("color is required")
	}
	
	if truck.Length <= 0 {
		return errors.New("length is required and must be positive")
	}
	
	if truck.Width <= 0 {
		return errors.New("width is required and must be positive")
	}
	
	if truck.Height <= 0 {
		return errors.New("height is required and must be positive")
	}
	
	if truck.MaxWeight <= 0 {
		return errors.New("max weight is required and must be positive")
	}
	
	// Validate truck type
	validTypes := map[string]bool{
		"Small":  true,
		"Medium": true,
		"Large":  true,
	}
	if !validTypes[truck.TruckType] {
		return errors.New("truck type must be one of: Small, Medium, Large")
	}
	
	return nil
}

// validatePhotoFile validates uploaded photo files
func (s *truckService) validatePhotoFile(file *multipart.FileHeader) error {
	// Check file size (max 10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxSize {
		return fmt.Errorf("file %s is too large (%d bytes), maximum size is 10MB", file.Filename, file.Size)
	}
	
	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	if !validExts[ext] {
		return fmt.Errorf("file %s has invalid extension %s, allowed: .jpg, .jpeg, .png, .gif", file.Filename, ext)
	}
	
	return nil
}