package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL      string
	MinIOURL         string
	MinIOUser        string
	MinIOPassword    string
}

func Load() (*Config, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	minioURL := os.Getenv("MINIO_URL")
	minioUser := os.Getenv("MINIO_ROOT_USER")
	minioPassword := os.Getenv("MINIO_ROOT_PASSWORD")

	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	if minioURL == "" || minioUser == "" || minioPassword == "" {
		return nil, fmt.Errorf("missing required MinIO environment variables")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	return &Config{
		DatabaseURL:   dbURL,
		MinIOURL:      minioURL,
		MinIOUser:     minioUser,
		MinIOPassword: minioPassword,
	}, nil
}
