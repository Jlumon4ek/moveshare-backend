package config

import (
	"log"
	"os"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

func LoadMinioConfig() *MinioConfig {
	endpoint := os.Getenv("MINIO_URL")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	bucket := "truck-photos"

	log.Printf("MINIO_URL: %q", endpoint)
	log.Printf("MINIO_ROOT_USER: %q", accessKey)
	log.Printf("MINIO_ROOT_PASSWORD: %q", secretKey)

	if endpoint == "" {
		log.Println("ERROR: MINIO_URL environment variable is not set!")
	}
	if accessKey == "" {
		log.Println("ERROR: MINIO_ROOT_USER environment variable is not set!")
	}
	if secretKey == "" {
		log.Println("ERROR: MINIO_ROOT_PASSWORD environment variable is not set!")
	}

	return &MinioConfig{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		UseSSL:    false,
		Bucket:    bucket,
	}
}
