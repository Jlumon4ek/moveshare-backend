package main

import (
	"context"
	"log"
	"moveshare/internal/config"
	"moveshare/internal/handlers"
	"moveshare/internal/handlers/chat"
	"moveshare/internal/handlers/review"
	"moveshare/internal/repository"
	chatRepo "moveshare/internal/repository/chat"
	"moveshare/internal/repository/company"
	"moveshare/internal/repository/payment"
	reviewRepo "moveshare/internal/repository/review"
	"moveshare/internal/repository/truck"
	"moveshare/internal/repository/user"
	"moveshare/internal/repository/verification"
	"moveshare/internal/repository/admin"
	"moveshare/internal/repository/password_reset"

	"moveshare/internal/router"
	"moveshare/internal/service"

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

	minioRepo, err := repository.MinioRepository(&cfg.Minio)
	if err != nil {
		log.Fatalf("failed to initialize Minio repository: %v", err)
	}
	minioService := service.NewMinioService(minioRepo)
	truckRepo := truck.NewTruckRepository(db)
	truckService := service.NewTruckService(truckRepo, minioRepo)

	verificationRepo := verification.NewVerificationRepository(db)
	verificationService := service.NewVerificationService(verificationRepo, minioRepo)

	stripeService := service.NewStripeService(&cfg.Stripe)

	paymentRepo := payment.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, stripeService, userService)

	// Password reset services
	passwordResetRepo := password_reset.NewPasswordResetRepository(db)
	emailService := service.NewEmailService()
	passwordResetService := service.NewPasswordResetService(passwordResetRepo, emailService)

	r := gin.Default()

	// Увеличиваем лимит для загрузки файлов до 100MB
	r.MaxMultipartMemory = 100 << 20 // 100 MB

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:5173",
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

	jobRepo := repository.NewJobRepository(db)
	jobService := service.NewJobService(jobRepo, &cfg.GoogleMaps, minioRepo)
	jobHandler := handlers.NewJobHandler(jobService, minioRepo)

	locationRepo := repository.NewLocationRepository(db)
	locationService := service.NewLocationService(locationRepo)
	locationHandler := handlers.NewLocationHandler(locationService)

	chatRepo := chatRepo.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo)

	reviewRepo := reviewRepo.NewReviewRepository(db)
	reviewService := service.NewReviewService(reviewRepo)
	reviewHandler := review.NewReviewHandler(reviewService)

	// Инициализация WebSocket hub
	hub := chat.NewHub(chatService)
	go hub.Run()

	// Добавь после routes.SetupLocationRoutes:

	apiGroup := r.Group("/api")
	{
		router.AdminRouter(apiGroup, jwtAuth, adminService)
		router.UserRouter(apiGroup, userService, minioService, jwtAuth, passwordResetService)
		router.CompanyRouter(apiGroup, companyService, jwtAuth)
		router.TruckRouter(apiGroup, truckService, jwtAuth)
		router.VerificationRouter(apiGroup, verificationService, jwtAuth)
		router.PaymentRouter(apiGroup, paymentService, jwtAuth) // ✅ Добавить
		router.SetupJobRoutes(apiGroup, jobHandler, jwtAuth)
		router.SetupLocationRoutes(apiGroup, locationHandler)
		router.SetupChatRoutes(apiGroup, chatService, *jobService, jwtAuth, hub)
		router.SetupReviewRoutes(apiGroup, reviewHandler, jwtAuth)
	}

	log.Println("Starting server on :8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
