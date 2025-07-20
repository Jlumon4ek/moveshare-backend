package main

import (
	"context"
	"log"
	_ "moveshare/docs" // Import generated docs for MoveShare API
	"moveshare/internal/auth"
	"moveshare/internal/config"
	"moveshare/internal/handlers"
	"moveshare/internal/repository"
	"moveshare/internal/service"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title MoveShare API
// @version 1.0
// @description API for user authentication, job management, and truck management in MoveShare application
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@moveshare.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT auth
	jwtAuth, err := auth.NewJWTAuth("keys/jwt-private.pem", "keys/jwt-public.pem")
	if err != nil {
		log.Fatalf("failed to initialize JWT auth: %v", err)
	}

	minioCfg := config.LoadMinioConfig()
	minioClient, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKey, minioCfg.SecretKey, ""),
		Secure: minioCfg.UseSSL,
	})
	if err != nil {
		log.Fatalf("failed to init minio: %v", err)
	}

	// Создайте bucket, если его нет
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, minioCfg.Bucket)
	if err != nil {
		log.Fatalf("failed to check minio bucket: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, minioCfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create minio bucket: %v", err)
		}
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	truckRepo := repository.NewTruckRepository(db) // New repository
	cardRepo := repository.NewCardRepository(db)   // New card repository

	// Initialize services
	userService := service.NewUserService(userRepo, jwtAuth)
	jobService := service.NewJobService(jobRepo)
	companyService := service.NewCompanyService(companyRepo)
	truckService := service.NewTruckService(truckRepo) // New service
	cardService := service.NewCardService(cardRepo)    // New card service

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	jobHandler := handlers.NewJobHandler(jobService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	truckHandler := handlers.NewTruckHandler(truckService, minioClient, minioCfg.Bucket)
	cardHandler := handlers.NewCardHandler(cardService) // New card handler

	// Setup router
	r := chi.NewRouter()

	// CORS configuration - allows all origins
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false, // Set to false when using "*" for AllowedOrigins
		MaxAge:           300,   // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes under /api/v1
	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Post("/sign-up", userHandler.SignUp)
		r.Post("/sign-in", userHandler.SignIn)

		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(jwtAuth))

			r.Route("/jobs", func(r chi.Router) {
				r.Post("/", jobHandler.CreateJob)
				r.Get("/available", jobHandler.GetAvailableJobs)
				r.Get("/my", jobHandler.GetUserJobs)
				r.Delete("/{id}", jobHandler.DeleteJob)
				r.Post("/{id}/apply", jobHandler.ApplyForJob)
				r.Get("/applications/my", jobHandler.GetMyApplications)
			})

			r.Route("/company", func(r chi.Router) {
				r.Get("/", companyHandler.GetCompany)
				r.Patch("/", companyHandler.PatchCompany)
			})

			r.Route("/trucks", func(r chi.Router) {
				r.Post("/", truckHandler.CreateTruck)
				r.Get("/", truckHandler.GetUserTrucks)
				r.Get("/{id}", truckHandler.GetTruckByID)
				r.Put("/{id}", truckHandler.UpdateTruck)
				r.Delete("/{id}", truckHandler.DeleteTruck)
			})

			r.Route("/cards", func(r chi.Router) {
				r.Post("/", cardHandler.CreateCard)
				r.Get("/", cardHandler.GetUserCards)
				r.Get("/{id}", cardHandler.GetCardByID)
				r.Put("/{id}", cardHandler.UpdateCard)
				r.Delete("/{id}", cardHandler.DeleteCard)
				r.Post("/{id}/default", cardHandler.SetDefaultCard)
			})
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler())

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
