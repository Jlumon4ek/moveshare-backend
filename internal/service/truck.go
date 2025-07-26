package service

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"moveshare/internal/repository"
	"moveshare/internal/repository/truck"
	"path/filepath"
	"time"
)

type TruckService interface {
	CreateTruck(ctx context.Context, truck *models.Truck) error
	DeleteTruck(ctx context.Context, id int64) error
	GetTruckByID(ctx context.Context, id int64) (*models.Truck, error)
	GetUserTrucks(ctx context.Context, userID int64) ([]*models.Truck, error)
	UpdateTruck(ctx context.Context, userID int64, truck *models.Truck) error
	InsertPhoto(ctx context.Context, truckID int64, objectID string) error
	GetTruckPhotos(ctx context.Context, truckID int64) ([]string, error)
	DeleteTruckPhoto(ctx context.Context, truckID int64, photoID string) error
}

type truckService struct {
	repo      truck.TruckRepository
	minioRepo *repository.Repository
}

func NewTruckService(repo truck.TruckRepository, minioRepo *repository.Repository) TruckService {
	return &truckService{repo: repo, minioRepo: minioRepo}
}

func (s *truckService) CreateTruck(ctx context.Context, truck *models.Truck) error {
	truck, err := s.repo.CreateTruck(ctx, truck)
	if err != nil {
		return fmt.Errorf("failed to create truck: %w", err)
	}

	for _, fileHeader := range truck.Photos {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open photo: %w", err)
		}
		defer file.Close()

		data := make([]byte, fileHeader.Size)
		_, err = file.Read(data)
		if err != nil {
			return fmt.Errorf("failed to read photo: %w", err)
		}
		ext := filepath.Ext(fileHeader.Filename)
		objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		err = s.minioRepo.UploadBytes(ctx, "trucks", objectName, data, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("failed to upload photo to MinIO: %w", err)
		}

		if err := s.repo.InsertPhoto(ctx, truck.ID, objectName); err != nil {
			return fmt.Errorf("failed to insert photo record: %w", err)
		}
	}

	return nil
}

func (s *truckService) DeleteTruck(ctx context.Context, id int64) error {
	if err := s.repo.DeleteTruck(ctx, id); err != nil {
		return fmt.Errorf("failed to delete truck: %w", err)
	}

	return nil
}

func (s *truckService) GetTruckByID(ctx context.Context, id int64) (*models.Truck, error) {
	truck, err := s.repo.GetTruckByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get truck by ID: %w", err)
	}

	photoIDs, err := s.repo.GetTruckPhotos(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get truck photos: %w", err)
	}

	var photoURLs []string
	for _, objectName := range photoIDs {
		url, err := s.minioRepo.GetFileURL(ctx, "trucks", objectName, 10*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("failed to generate URL for %s: %w", objectName, err)
		}
		photoURLs = append(photoURLs, url)
	}
	truck.PhotoURLs = photoURLs
	return truck, nil

}

func (s *truckService) GetUserTrucks(ctx context.Context, userID int64) ([]*models.Truck, error) {
	trucks, err := s.repo.GetUserTrucks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user trucks: %w", err)
	}

	for _, truck := range trucks {
		photoIDs, err := s.repo.GetTruckPhotos(ctx, truck.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get photos for truck %d: %w", truck.ID, err)
		}

		var urls []string
		for _, objectName := range photoIDs {
			url, err := s.minioRepo.GetFileURL(ctx, "trucks", objectName, 10*time.Minute)
			if err != nil {
				return nil, fmt.Errorf("failed to get URL for %s: %w", objectName, err)
			}
			urls = append(urls, url)
		}
		truck.PhotoURLs = urls
	}

	return trucks, nil
}

func (s *truckService) UpdateTruck(ctx context.Context, userID int64, truck *models.Truck) error {
	return s.repo.UpdateTruck(ctx, userID, truck)
}

func (s *truckService) InsertPhoto(ctx context.Context, truckID int64, objectID string) error {
	return s.repo.InsertPhoto(ctx, truckID, objectID)
}

func (s *truckService) GetTruckPhotos(ctx context.Context, truckID int64) ([]string, error) {
	return s.repo.GetTruckPhotos(ctx, truckID)
}
func (s *truckService) DeleteTruckPhoto(ctx context.Context, truckID int64, photoID string) error {
	return s.repo.DeleteTruckPhoto(ctx, truckID, photoID)
}
