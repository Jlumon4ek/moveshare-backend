package service

import (
	"context"
	"errors"
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
	// Truck operations
	CreateTruck(ctx context.Context, userID int64, truck *repository.Truck) error
	GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error)
	GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error)
	UpdateTruck(ctx context.Context, userID, truckID int64, truck *repository.Truck) error
	DeleteTruck(ctx context.Context, userID, truckID int64) error
	
	// Photo operations
	UploadTruckPhotos(ctx context.Context, userID, truckID int64, files []*multipart.FileHeader) ([]repository.TruckPhoto, error)
	GetTruckPhotos(ctx context.Context, userID, truckID int64) ([]repository.TruckPhoto, error)
	DeleteTruckPhoto(ctx context.Context, userID, truckID, photoID int64) error
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

// GetUserTrucks fetches all trucks belonging to a user
func (s *truckService) GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error) {
	return s.truckRepo.GetUserTrucks(ctx, userID)
}

// GetTruckByID fetches a specific truck by ID for a user
func (s *truckService) GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error) {
	truck, err := s.truckRepo.GetTruckByID(ctx, userID, truckID)
	if err != nil {
		return nil, err
	}
	if truck == nil {
		return nil, errors.New("truck not found")
	}
	return truck, nil
}

// UpdateTruck updates an existing truck with validation
func (s *truckService) UpdateTruck(ctx context.Context, userID, truckID int64, truck *repository.Truck) error {
	if err := s.validateTruck(truck); err != nil {
		return err
	}
	
	// Verify the truck exists and belongs to the user
	existingTruck, err := s.truckRepo.GetTruckByID(ctx, userID, truckID)
	if err != nil {
		return err
	}
	if existingTruck == nil {
		return errors.New("truck not found")
	}
	
	truck.UserID = userID
	return s.truckRepo.UpdateTruck(ctx, userID, truckID, truck)
}

// DeleteTruck deletes a truck
func (s *truckService) DeleteTruck(ctx context.Context, userID, truckID int64) error {
	// Verify the truck exists and belongs to the user
	existingTruck, err := s.truckRepo.GetTruckByID(ctx, userID, truckID)
	if err != nil {
		return err
	}
	if existingTruck == nil {
		return errors.New("truck not found")
	}
	
	return s.truckRepo.DeleteTruck(ctx, userID, truckID)
}

// UploadTruckPhotos handles uploading multiple photos for a truck
func (s *truckService) UploadTruckPhotos(ctx context.Context, userID, truckID int64, files []*multipart.FileHeader) ([]repository.TruckPhoto, error) {
	// Verify the truck exists and belongs to the user
	truck, err := s.truckRepo.GetTruckByID(ctx, userID, truckID)
	if err != nil {
		return nil, err
	}
	if truck == nil {
		return nil, errors.New("truck not found")
	}
	
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}
	
	// Create upload directory if it doesn't exist
	uploadDir := "uploads/trucks"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}
	
	var uploadedPhotos []repository.TruckPhoto
	
	for _, fileHeader := range files {
		// Validate file
		if err := s.validatePhotoFile(fileHeader); err != nil {
			return nil, err
		}
		
		// Generate unique filename
		filename := s.generatePhotoFilename(userID, truckID, fileHeader.Filename)
		filePath := filepath.Join(uploadDir, filename)
		
		// Save file to disk
		if err := s.saveUploadedFile(fileHeader, filePath); err != nil {
			return nil, fmt.Errorf("failed to save file %s: %w", fileHeader.Filename, err)
		}
		
		// Create database record
		photo := &repository.TruckPhoto{
			TruckID:  truckID,
			UserID:   userID,
			FileName: fileHeader.Filename,
			FileURL:  filePath,
		}
		
		if err := s.truckRepo.AddTruckPhoto(ctx, photo); err != nil {
			// Clean up the saved file if database operation fails
			os.Remove(filePath)
			return nil, fmt.Errorf("failed to save photo record: %w", err)
		}
		
		uploadedPhotos = append(uploadedPhotos, *photo)
	}
	
	return uploadedPhotos, nil
}

// GetTruckPhotos fetches all photos for a truck
func (s *truckService) GetTruckPhotos(ctx context.Context, userID, truckID int64) ([]repository.TruckPhoto, error) {
	return s.truckRepo.GetTruckPhotos(ctx, userID, truckID)
}

// DeleteTruckPhoto deletes a specific truck photo
func (s *truckService) DeleteTruckPhoto(ctx context.Context, userID, truckID, photoID int64) error {
	// Get the photo details first to get the file path
	photos, err := s.truckRepo.GetTruckPhotos(ctx, userID, truckID)
	if err != nil {
		return err
	}
	
	var photoToDelete *repository.TruckPhoto
	for _, photo := range photos {
		if photo.ID == photoID {
			photoToDelete = &photo
			break
		}
	}
	
	if photoToDelete == nil {
		return errors.New("photo not found")
	}
	
	// Delete from database first
	if err := s.truckRepo.DeleteTruckPhoto(ctx, userID, truckID, photoID); err != nil {
		return err
	}
	
	// Delete the physical file (ignore errors as the file might not exist)
	os.Remove(photoToDelete.FileURL)
	
	return nil
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
	if truck.Year < 1900 || truck.Year > time.Now().Year()+1 {
		return errors.New("invalid year")
	}
	if truck.TruckType != repository.TruckTypeSmall && 
	   truck.TruckType != repository.TruckTypeMedium && 
	   truck.TruckType != repository.TruckTypeLarge {
		return errors.New("invalid truck type, must be Small, Medium, or Large")
	}
	if truck.Length < 0 || truck.Width < 0 || truck.Height < 0 || truck.MaxWeight < 0 {
		return errors.New("dimensions and max weight must be positive")
	}
	
	return nil
}

// validatePhotoFile validates uploaded photo files
func (s *truckService) validatePhotoFile(fileHeader *multipart.FileHeader) error {
	// Check file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file %s is too large (max 10MB)", fileHeader.Filename)
	}
	
	// Check file extension
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return nil
		}
	}
	
	return fmt.Errorf("file %s has invalid extension, allowed: %v", fileHeader.Filename, allowedExtensions)
}

// generatePhotoFilename generates a unique filename for uploaded photos
func (s *truckService) generatePhotoFilename(userID, truckID int64, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("truck_%d_user_%d_%d%s", truckID, userID, timestamp, ext)
}

// saveUploadedFile saves a multipart file to disk
func (s *truckService) saveUploadedFile(fileHeader *multipart.FileHeader, filePath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, src)
	return err
}