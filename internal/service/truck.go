package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"moveshare/internal/config"
	"moveshare/internal/repository"
)

// TruckService defines the interface for truck business logic
type TruckService interface {
	CreateTruck(ctx context.Context, truck *repository.Truck, photos []*multipart.FileHeader) error
	GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error)
	GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error)
	UpdateTruck(ctx context.Context, truck *repository.Truck) error
	DeleteTruck(ctx context.Context, userID, truckID int64) error
}

// truckService implements TruckService
type truckService struct {
	truckRepo   repository.TruckRepository
	minioClient *minio.Client
	cfg         *config.Config
}

// NewTruckService creates a new TruckService
func NewTruckService(truckRepo repository.TruckRepository, cfg *config.Config) (TruckService, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.MinIOURL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOUser, cfg.MinIOPassword, ""),
		Secure: false, // Set to true for HTTPS
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	return &truckService{
		truckRepo:   truckRepo,
		minioClient: minioClient,
		cfg:         cfg,
	}, nil
}

// CreateTruck creates a new truck with photos
func (s *truckService) CreateTruck(ctx context.Context, truck *repository.Truck, photos []*multipart.FileHeader) error {
	// Create truck in database first
	err := s.truckRepo.CreateTruck(ctx, truck)
	if err != nil {
		return fmt.Errorf("failed to create truck: %w", err)
	}

	// Upload photos if provided
	if len(photos) > 0 {
		err = s.uploadTruckPhotos(ctx, truck.ID, truck.UserID, photos)
		if err != nil {
			// If photo upload fails, we should probably delete the truck
			// but for now just return the error
			return fmt.Errorf("failed to upload photos: %w", err)
		}
	}

	return nil
}

// GetUserTrucks fetches all trucks for a user
func (s *truckService) GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error) {
	return s.truckRepo.GetUserTrucks(ctx, userID)
}

// GetTruckByID fetches a specific truck by ID
func (s *truckService) GetTruckByID(ctx context.Context, userID, truckID int64) (*repository.Truck, error) {
	return s.truckRepo.GetTruckByID(ctx, userID, truckID)
}

// UpdateTruck updates an existing truck
func (s *truckService) UpdateTruck(ctx context.Context, truck *repository.Truck) error {
	return s.truckRepo.UpdateTruck(ctx, truck)
}

// DeleteTruck deletes a truck and its photos
func (s *truckService) DeleteTruck(ctx context.Context, userID, truckID int64) error {
	// Get truck photos to delete from MinIO
	photos, err := s.truckRepo.GetTruckPhotos(ctx, truckID)
	if err != nil {
		return fmt.Errorf("failed to get truck photos: %w", err)
	}

	// Delete truck from database (this will also delete photo records)
	err = s.truckRepo.DeleteTruck(ctx, userID, truckID)
	if err != nil {
		return fmt.Errorf("failed to delete truck: %w", err)
	}

	// Delete photos from MinIO
	bucketName := fmt.Sprintf("trucks-user-%d", userID)
	for _, photo := range photos {
		err = s.minioClient.RemoveObject(ctx, bucketName, photo.FileName, minio.RemoveObjectOptions{})
		if err != nil {
			// Log error but don't fail the entire operation
			fmt.Printf("Failed to delete photo %s from MinIO: %v\n", photo.FileName, err)
		}
	}

	return nil
}

// uploadTruckPhotos uploads photos to MinIO and saves records to database
func (s *truckService) uploadTruckPhotos(ctx context.Context, truckID, userID int64, photos []*multipart.FileHeader) error {
	bucketName := fmt.Sprintf("trucks-user-%d", userID)

	// Ensure bucket exists
	err := s.ensureBucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	for _, fileHeader := range photos {
		// Validate file type
		if !s.isValidImageFile(fileHeader.Filename) {
			return fmt.Errorf("invalid file type: %s", fileHeader.Filename)
		}

		// Generate unique filename
		ext := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprintf("truck-%d-%s%s", truckID, uuid.New().String(), ext)

		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		// Upload to MinIO
		_, err = s.minioClient.PutObject(ctx, bucketName, fileName, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: getContentType(ext),
		})
		if err != nil {
			return fmt.Errorf("failed to upload file to MinIO: %w", err)
		}

		// Generate URL for the uploaded file
		fileURL := fmt.Sprintf("http://%s/%s/%s", s.cfg.MinIOURL, bucketName, fileName)

		// Save photo record in database
		photo := &repository.TruckPhoto{
			TruckID:  truckID,
			UserID:   userID,
			FileName: fileName,
			FileURL:  fileURL,
		}

		err = s.truckRepo.AddTruckPhoto(ctx, photo)
		if err != nil {
			return fmt.Errorf("failed to save photo record: %w", err)
		}
	}

	return nil
}

// ensureBucketExists creates bucket if it doesn't exist
func (s *truckService) ensureBucketExists(ctx context.Context, bucketName string) error {
	exists, err := s.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = s.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// isValidImageFile checks if the file is a valid image
func (s *truckService) isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	
	return false
}

// getContentType returns the MIME type for a file extension
func getContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}