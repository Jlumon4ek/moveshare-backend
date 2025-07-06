package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"moveshare/internal/repository"
)

// TruckService defines the interface for truck business logic
type TruckService interface {
	CreateTruck(ctx context.Context, userID int64, truck *repository.Truck) error
	GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error)
	DeleteTruck(ctx context.Context, userID, truckID int64) error
	UploadTruckPhoto(ctx context.Context, userID, truckID int64, file multipart.File, header *multipart.FileHeader) (*repository.TruckPhoto, error)
	GetTruckPhotos(ctx context.Context, truckID int64) ([]repository.TruckPhoto, error)
}

// truckService implements TruckService
type truckService struct {
	truckRepo repository.TruckRepository
}

// NewTruckService creates a new TruckService
func NewTruckService(truckRepo repository.TruckRepository) TruckService {
	return &truckService{truckRepo: truckRepo}
}

// CreateTruck creates a new truck
func (s *truckService) CreateTruck(ctx context.Context, userID int64, truck *repository.Truck) error {
	// Set user ID
	truck.UserID = userID
	
	// Validate required fields
	if truck.TruckName == "" {
		return fmt.Errorf("truck name is required")
	}
	if truck.LicensePlate == "" {
		return fmt.Errorf("license plate is required")
	}
	if truck.Make == "" {
		return fmt.Errorf("make is required")
	}
	if truck.Model == "" {
		return fmt.Errorf("model is required")
	}
	if truck.Year == 0 {
		return fmt.Errorf("year is required")
	}
	if truck.Color == "" {
		return fmt.Errorf("color is required")
	}
	
	// Validate truck type
	validTypes := []string{"Small", "Medium", "Large"}
	valid := false
	for _, t := range validTypes {
		if truck.TruckType == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid truck type. Must be one of: Small, Medium, Large")
	}
	
	return s.truckRepo.CreateTruck(ctx, truck)
}

// GetUserTrucks fetches all trucks for a user
func (s *truckService) GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error) {
	return s.truckRepo.GetUserTrucks(ctx, userID)
}

// DeleteTruck deletes a truck owned by the user
func (s *truckService) DeleteTruck(ctx context.Context, userID, truckID int64) error {
	return s.truckRepo.DeleteTruck(ctx, userID, truckID)
}

// UploadTruckPhoto handles photo upload for a truck
func (s *truckService) UploadTruckPhoto(ctx context.Context, userID, truckID int64, file multipart.File, header *multipart.FileHeader) (*repository.TruckPhoto, error) {
	// Verify truck belongs to user
	truck, err := s.truckRepo.GetTruckByID(ctx, truckID)
	if err != nil {
		return nil, fmt.Errorf("truck not found")
	}
	if truck.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to truck")
	}
	
	// Validate file type
	filename := header.Filename
	if !isValidImageFile(filename) {
		return nil, fmt.Errorf("invalid file type. Only jpg, jpeg, png files are allowed")
	}
	
	// Create upload directory if it doesn't exist
	uploadDir := "uploads/trucks"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}
	
	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("truck_%d_%d%s", truckID, timestamp, ext)
	filePath := filepath.Join(uploadDir, uniqueFilename)
	
	// Save file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()
	
	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}
	
	// Create photo record in database
	photo := &repository.TruckPhoto{
		TruckID:  truckID,
		UserID:   userID,
		FileName: uniqueFilename,
		FileURL:  "/uploads/trucks/" + uniqueFilename,
	}
	
	if err := s.truckRepo.UploadTruckPhoto(ctx, photo); err != nil {
		// Clean up file if database operation fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save photo record: %v", err)
	}
	
	return photo, nil
}

// GetTruckPhotos fetches all photos for a truck
func (s *truckService) GetTruckPhotos(ctx context.Context, truckID int64) ([]repository.TruckPhoto, error) {
	return s.truckRepo.GetTruckPhotos(ctx, truckID)
}

// isValidImageFile checks if the file is a valid image type
func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}