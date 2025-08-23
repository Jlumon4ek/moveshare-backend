package main

import (
	"context"
	"log"
	"moveshare/internal/config"
	"moveshare/internal/handlers"
	"moveshare/internal/handlers/chat"
	"moveshare/internal/handlers/review"
	"moveshare/internal/repository"
	"moveshare/internal/repository/admin"
	chatRepo "moveshare/internal/repository/chat"
	"moveshare/internal/repository/company"
	notificationRepo "moveshare/internal/repository/notifications"
	"moveshare/internal/repository/password_reset"
	"moveshare/internal/repository/payment"
	reviewRepo "moveshare/internal/repository/review"
	sessionRepo "moveshare/internal/repository/session"
	"moveshare/internal/repository/truck"
	"moveshare/internal/repository/user"
	"moveshare/internal/repository/verification"
	"moveshare/internal/websocket"
	"strings"

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

	sessionRepository := sessionRepo.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepository)

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

	// Добавляем middleware для детального логирования запросов
	r.Use(func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "upload") {
			log.Printf("=== UPLOAD REQUEST START ===")
			log.Printf("Method: %s", c.Request.Method)
			log.Printf("URL: %s", c.Request.URL.String())
			log.Printf("Content-Type: %s", c.Request.Header.Get("Content-Type"))
			log.Printf("Content-Length: %s", c.Request.Header.Get("Content-Length"))
			log.Printf("Remote Addr: %s", c.Request.RemoteAddr)
			log.Printf("User-Agent: %s", c.Request.Header.Get("User-Agent"))

			// Логируем размер тела запроса
			if c.Request.ContentLength > 0 {
				log.Printf("Request body size: %d bytes (%.2f MB)", c.Request.ContentLength, float64(c.Request.ContentLength)/(1024*1024))
			}
		}
		c.Next()
		if strings.Contains(c.Request.URL.Path, "upload") {
			log.Printf("=== UPLOAD REQUEST END ===")
		}
	})

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

	// Сервис уведомлений (инициализируем раньше)
	notificationHub := websocket.NewNotificationHub()
	go notificationHub.Run()
	notificationRepoInstance := notificationRepo.NewNotificationRepository(db)
	notificationService := service.NewNotificationService(notificationHub, notificationRepoInstance)

	jobService := service.NewJobService(jobRepo, &cfg.GoogleMaps, minioRepo, notificationService)

	locationRepo := repository.NewLocationRepository(db)
	locationService := service.NewLocationService(locationRepo)
	locationHandler := handlers.NewLocationHandler(locationService)

	chatRepo := chatRepo.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo)

	jobHandler := handlers.NewJobHandler(jobService, chatService, notificationService, minioRepo, paymentService, adminService)

	reviewRepo := reviewRepo.NewReviewRepository(db)
	reviewService := service.NewReviewService(reviewRepo)
	reviewHandler := review.NewReviewHandler(reviewService, notificationService, jobService)

	// Инициализация WebSocket hub для чата
	hub := chat.NewHub(chatService)
	go hub.Run()

	// Добавь после routes.SetupLocationRoutes:

	apiGroup := r.Group("/api")
	{
		router.AdminRouter(apiGroup, jwtAuth, adminService)
		router.UserRouter(apiGroup, userService, minioService, jwtAuth, passwordResetService, sessionService)
		router.CompanyRouter(apiGroup, companyService, jwtAuth)
		router.TruckRouter(apiGroup, truckService, jwtAuth)
		router.VerificationRouter(apiGroup, verificationService, jwtAuth)
		router.PaymentRouter(apiGroup, paymentService, jwtAuth)
		router.SetupJobRoutes(apiGroup, jobHandler, jwtAuth)
		router.SetupLocationRoutes(apiGroup, locationHandler)
		router.SetupChatRoutes(apiGroup, chatService, *jobService, jwtAuth, hub, notificationService)
		router.SetupNotificationRoutes(apiGroup, jwtAuth, notificationHub, notificationService)
		router.SetupReviewRoutes(apiGroup, reviewHandler, jwtAuth)

		// Public system settings route
		apiGroup.GET("/commission-rate", handlers.GetCommissionRate(adminService))
	}

	log.Println("Starting server on :8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
