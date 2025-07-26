package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"moveshare/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Repository struct {
	client *minio.Client
	config *config.MinioConfig
}

func MinioRepository(cfg *config.MinioConfig) (*Repository, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	return &Repository{
		client: client,
		config: cfg,
	}, nil
}

func (r *Repository) ensureBucketExists(ctx context.Context, bucketName string) error {
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Bucket %q created\n", bucketName)
	} else {
		log.Printf("Bucket %q already exists\n", bucketName)
	}
	return nil
}

func (r *Repository) UploadFile(ctx context.Context, bucket, objectName, filePath string) error {
	if err := r.ensureBucketExists(ctx, bucket); err != nil {
		return err
	}
	_, err := r.client.FPutObject(ctx, bucket, objectName, filePath, minio.PutObjectOptions{})
	return err
}

func (r *Repository) UploadBytes(ctx context.Context, bucket, objectName string, data []byte, contentType string) error {
	if err := r.ensureBucketExists(ctx, bucket); err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	_, err := r.client.PutObject(ctx, bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (r *Repository) DownloadFile(ctx context.Context, bucket, objectName, filePath string) error {
	return r.client.FGetObject(ctx, bucket, objectName, filePath, minio.GetObjectOptions{})
}

func (r *Repository) GetFileURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error) {
	url, err := r.client.PresignedGetObject(ctx, bucket, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (r *Repository) ListObjects(ctx context.Context, bucket, prefix string) ([]string, error) {
	if err := r.ensureBucketExists(ctx, bucket); err != nil {
		return nil, err
	}

	var result []string
	for object := range r.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if object.Err != nil {
			return nil, object.Err
		}
		result = append(result, object.Key)
	}
	return result, nil
}

func (r *Repository) DeleteObject(ctx context.Context, bucket, objectName string) error {
	return r.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

func (r *Repository) DownloadStream(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	return r.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
}
