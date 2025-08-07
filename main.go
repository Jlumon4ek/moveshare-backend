package main

import (
	"context"
	"log"
	"moveshare/internal/config"
	"moveshare/internal/repository"
	"moveshare/internal/repository/admin"
	"moveshare/internal/repository/chat"
	"moveshare/internal/repository/company"
	"moveshare/internal/repository/job"
	"moveshare/internal/repository/truck"
	"moveshare/internal/repository/user"
	"moveshare/internal/repository/verification"

	"moveshare/internal/router"
	"moveshare/internal/service"

	chathandlers "moveshare/internal/handlers/chat"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "moveshare/docs"
)

// @title MoveShare API
// @version 1.0
// @description API для приложения MoveShare
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter JWT token in format: Bearer {your_token_here}
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

	jwtAuth, err := service.NewJWTAuth("keys/jwt-private.pem", "keys/jwt-public.pem")
	if err != nil {
		log.Fatalf("failed to initialize JWT auth: %v", err)
	}

	adminRepo := admin.NewAdminRepository(db)
	adminService := service.NewAdminService(adminRepo)

	userRepo := user.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	companyRepo := company.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepo, userRepo)

	jobRepo := job.NewJobRepository(db)
	jobService := service.NewJobService(jobRepo)

	minioRepo, err := repository.MinioRepository(&cfg.Minio)
	if err != nil {
		log.Fatalf("failed to initialize Minio repository: %v", err)
	}
	truckRepo := truck.NewTruckRepository(db)
	truckService := service.NewTruckService(truckRepo, minioRepo)

	verificationRepo := verification.NewVerificationRepository(db)
	verificationService := service.NewVerificationService(verificationRepo, minioRepo)

	chatRepo := chat.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo)

	chatHub := chathandlers.NewHub(chatService)
	go chatHub.Run()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:5173", // Vite
		"http://127.0.0.1:5173",
		"http://127.0.0.1:3000",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}
	config.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := r.Group("/api")
	{
		router.AdminRouter(apiGroup, jwtAuth, adminService)
		router.UserRouter(apiGroup, userService, jwtAuth)
		router.CompanyRouter(apiGroup, companyService, jwtAuth)
		router.JobRouter(apiGroup, jobService, jwtAuth)
		router.TruckRouter(apiGroup, truckService, jwtAuth)
		router.VerificationRouter(apiGroup, verificationService, jwtAuth)
		router.ChatRouter(apiGroup, chatService, jwtAuth, chatHub)
	}

	log.Println("Starting server on :8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
