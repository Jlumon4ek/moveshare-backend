package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "moveshare/docs" // Import generated docs for MoveShare API

	"moveshare/internal/auth"
	"moveshare/internal/config"
	"moveshare/internal/handlers"
	"moveshare/internal/repository"
	"moveshare/internal/service"

	"github.com/go-chi/cors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	// Initialize services
	userService := service.NewUserService(userRepo, jwtAuth)
	jobService := service.NewJobService(jobRepo)
	companyService := service.NewCompanyService(companyRepo)
	truckService := service.NewTruckService(truckRepo) // New service

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	jobHandler := handlers.NewJobHandler(jobService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	truckHandler := handlers.NewTruckHandler(truckService, minioClient, minioCfg.Bucket)

	// Setup router
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
		AllowedHeaders: []string{"*"},
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files for uploaded photos
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads/"))))

	// API routes under /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/sign-up", userHandler.SignUp)
		r.Post("/sign-in", userHandler.SignIn)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(jwtAuth))

			// Job routes
			r.Post("/jobs", jobHandler.CreateJob)
			r.Get("/jobs/available", jobHandler.GetAvailableJobs)
			r.Get("/jobs/my", jobHandler.GetUserJobs)
			r.Delete("/jobs/{id}", jobHandler.DeleteJob)
			r.Post("/jobs/{id}/apply", jobHandler.ApplyForJob)
			r.Get("/jobs/applications/my", jobHandler.GetMyApplications)

			// Company routes
			r.Get("/company", companyHandler.GetCompany)
			r.Patch("/company", companyHandler.PatchCompany)

			// Truck routes
			r.Post("/trucks", truckHandler.CreateTruck)
			r.Get("/trucks", truckHandler.GetUserTrucks)
			r.Get("/trucks/{id}", truckHandler.GetTruckByID)
			r.Put("/trucks/{id}", truckHandler.UpdateTruck)
			r.Delete("/trucks/{id}", truckHandler.DeleteTruck)
		})
	})

	// Swagger documentation
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
