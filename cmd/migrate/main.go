package main

import (
	"context"
	"log"
	"moveshare/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Apply migration for city/state fields
	migrations := []string{
		"ALTER TABLE jobs ADD COLUMN IF NOT EXISTS pickup_city VARCHAR(100) DEFAULT '';",
		"ALTER TABLE jobs ADD COLUMN IF NOT EXISTS pickup_state VARCHAR(50) DEFAULT '';",
		"ALTER TABLE jobs ADD COLUMN IF NOT EXISTS delivery_city VARCHAR(100) DEFAULT '';", 
		"ALTER TABLE jobs ADD COLUMN IF NOT EXISTS delivery_state VARCHAR(50) DEFAULT '';",
	}

	for _, migration := range migrations {
		_, err := db.Exec(context.Background(), migration)
		if err != nil {
			log.Fatalf("Failed to execute migration: %s\nError: %v", migration, err)
		}
		log.Printf("Successfully executed: %s", migration)
	}

	log.Println("All migrations applied successfully!")
}