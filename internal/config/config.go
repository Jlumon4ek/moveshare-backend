package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Minio    MinioConfig
	Stripe   StripeConfig
	Email    EmailConfig
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type DatabaseConfig struct {
	URL string
}

type StripeConfig struct {
	PublicKey     string
	PrivateKey    string
	RestrictedKey string
}

type EmailConfig struct {
	SendPulseID     string
	SendPulseSecret string
}

func Load() (*Config, error) {
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	minioConfig := loadMinioConfig()

	return &Config{
		Database: *dbConfig,
		Minio:    *minioConfig,
	}, nil
}

func loadDatabaseConfig() (*DatabaseConfig, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	return &DatabaseConfig{
		URL: dbURL,
	}, nil
}

func loadMinioConfig() *MinioConfig {
	endpoint := os.Getenv("MINIO_URL")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	bucket := "truck-photos"

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

func loadStripeConfig() (*StripeConfig, error) {
	publicKey := os.Getenv("STRIPE_PUBLIC_KEY")
	privateKey := os.Getenv("STRIPE_PRIVATE_KEY")
	restrictedKey := os.Getenv("STRIPE_RESTRICTED_KEY")

	if publicKey == "" || privateKey == "" || restrictedKey == "" {
		return nil, fmt.Errorf("missing required Stripe environment variables")
	}

	return &StripeConfig{
		PublicKey:     publicKey,
		PrivateKey:    privateKey,
		RestrictedKey: restrictedKey,
	}, nil
}

func loadEmailConfig() (*EmailConfig, error) {
	sendPulseID := os.Getenv("SENDPULSE_ID")
	sendPulseSecret := os.Getenv("SENDPULSE_SECRET_KEY")

	if sendPulseID == "" || sendPulseSecret == "" {
		return nil, fmt.Errorf("missing required SendPulse environment variables")
	}
	return &EmailConfig{
		SendPulseID:     sendPulseID,
		SendPulseSecret: sendPulseSecret,
	}, nil
}
