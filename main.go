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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	truckRepo := repository.NewTruckRepository(db)
	
	// Initialize services
	userService := service.NewUserService(userRepo, jwtAuth)
	jobService := service.NewJobService(jobRepo)
	companyService := service.NewCompanyService(companyRepo)
	truckService, err := service.NewTruckService(truckRepo, cfg)
	if err != nil {
		log.Fatalf("failed to initialize truck service: %v", err)
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	jobHandler := handlers.NewJobHandler(jobService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	truckHandler := handlers.NewTruckHandler(truckService)
	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes under /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/sign-up", userHandler.SignUp)
		r.Post("/sign-in", userHandler.SignIn)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(jwtAuth))
			r.Post("/jobs", jobHandler.CreateJob)
			r.Get("/jobs/available", jobHandler.GetAvailableJobs)
			r.Get("/jobs/my", jobHandler.GetUserJobs)
			r.Delete("/jobs/{id}", jobHandler.DeleteJob)
			r.Post("/jobs/{id}/apply", jobHandler.ApplyForJob)
			r.Get("/jobs/applications/my", jobHandler.GetMyApplications)
			r.Get("/company", companyHandler.GetCompany)
			r.Patch("/company", companyHandler.PatchCompany)
			
			// Truck routes
			r.Get("/trucks", truckHandler.GetUserTrucks)
			r.Post("/trucks", truckHandler.CreateTruck)
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
