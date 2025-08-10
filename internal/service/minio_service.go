package service

import (
	"context"
	"time"

	"moveshare/internal/repository"
)

type MinioService struct {
	minioRepo *repository.Repository
}

func NewMinioService(minioRepo *repository.Repository) *MinioService {
	return &MinioService{
		minioRepo: minioRepo,
	}
}

func (s *MinioService) UploadProfilePhoto(ctx context.Context, objectID string, data []byte, contentType string) error {
	bucket := "profile-photos"
	return s.minioRepo.UploadBytes(ctx, bucket, objectID, data, contentType)
}

func (s *MinioService) GetProfilePhotoURL(objectID string, expires time.Duration) (string, error) {
	bucket := "profile-photos"
	return s.minioRepo.GetFileURL(context.Background(), bucket, objectID, expires)
}

func (s *MinioService) DeleteProfilePhoto(objectID string) error {
	bucket := "profile-photos"
	return s.minioRepo.DeleteObject(context.Background(), bucket, objectID)
}