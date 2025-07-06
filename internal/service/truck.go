package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"moveshare/internal/repository"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// TruckService defines the interface for truck business logic
type TruckService interface {
	GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error)
	CreateTruck(ctx context.Context, userID int64, truck *repository.Truck, files []multipart.FileHeader) error
}

// truckService implements TruckService
type truckService struct {
	truckRepo   repository.TruckRepository
	minioClient *minio.Client
	bucketName  string
}

// NewTruckService creates a new TruckService
func NewTruckService(truckRepo repository.TruckRepository) TruckService {
	// Initialize Minio client with environment variables
	endpoint := "minio:9000"
	accessKey := "minioadmin"
	secretKey := "minioadmin123"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Minio client: %v", err)
	}

	// Create bucket if it doesn't exist
	bucketName := "truck-photos"
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists", bucketName)
		} else {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}

	return &truckService{
		truckRepo:   truckRepo,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

// GetUserTrucks fetches all trucks for the given userID
func (s *truckService) GetUserTrucks(ctx context.Context, userID int64) ([]repository.Truck, error) {
	trucks, err := s.truckRepo.GetUserTrucks(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Generate pre-signed URLs for photos
	for i := range trucks {
		for j, url := range trucks[i].Photos {
			reqParams := make(url.Values)
			presignedURL, err := s.minioClient.PresignedGetObject(context.Background(), s.bucketName, filepath.Base(url), time.Hour*24, reqParams)
			if err != nil {
				return nil, err
			}
			trucks[i].Photos[j] = presignedURL.String()
		}
	}

	return trucks, nil
}

// CreateTruck creates a new truck and uploads photos to Minio
func (s *truckService) CreateTruck(ctx context.Context, userID int64, truck *repository.Truck, files []multipart.FileHeader) error {
	if truck.TruckName == "" || truck.LicensePlate == "" {
		return errors.New("truck name and license plate are required")
	}
	truck.UserID = userID

	var photoURLs []string
	for _, file := range files {
		// Generate a unique file name
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%d-%d%s", userID, time.Now().UnixNano(), ext)
		contentType := file.Header.Get("Content-Type")

		// Open the file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Upload to Minio
		_, err = s.minioClient.PutObject(ctx, s.bucketName, fileName, src, file.Size, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			return err
		}
		// Store the internal URL for database
		photoURLs = append(photoURLs, fileName)
	}

	// Create truck and get it back with photos
	createdTruck, err := s.truckRepo.CreateTruck(ctx, truck, photoURLs)
	if err != nil {
		return err
	}
	*truck = *createdTruck // Update the input truck with the created data

	return nil
}
