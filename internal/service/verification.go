package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"moveshare/internal/repository"
	"moveshare/internal/repository/verification"
	"path/filepath"
	"time"
)

type VerificationService interface {
	InsertFileID(ctx context.Context, file *multipart.FileHeader, userID int64, fileType string) error
}

type verificationService struct {
	verificationRepo verification.VerificationRepository
	minioRepo        *repository.Repository
}

func NewVerificationService(verificationRepo verification.VerificationRepository, minioRepo *repository.Repository) VerificationService {
	return &verificationService{
		verificationRepo: verificationRepo,
		minioRepo:        minioRepo,
	}
}

func (s *verificationService) InsertFileID(ctx context.Context, file *multipart.FileHeader, userID int64, fileType string) error {
	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	data := make([]byte, file.Size)
	_, err = fileReader.Read(data)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	err = s.minioRepo.UploadBytes(ctx, "verification", objectName, data, file.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return s.verificationRepo.InsertFileID(ctx, userID, objectName, fileType)
}
