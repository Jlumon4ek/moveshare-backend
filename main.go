package main

import (
	"context"
	"log"
	"moveshare/internal/config"
	"moveshare/internal/repository"
	"moveshare/internal/repository/admin"
	"moveshare/internal/repository/company"
	"moveshare/internal/repository/job"
	"moveshare/internal/repository/truck"
	"moveshare/internal/repository/user"

	"moveshare/internal/router"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "moveshare/docs"
)

// @title MoveShare API
// @version 1.0
// @description API для приложения MoveShare

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
	companyService := service.NewCompanyService(companyRepo)

	jobRepo := job.NewJobRepository(db)
	jobService := service.NewJobService(jobRepo)

	minioRepo, err := repository.MinioRepository(&cfg.Minio)
	if err != nil {
		log.Fatalf("failed to initialize Minio repository: %v", err)
	}
	truckRepo := truck.NewTruckRepository(db)
	truckService := service.NewTruckService(truckRepo, minioRepo)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.AdminRouter(r, jwtAuth, adminService)
	router.UserRouter(r, userService, jwtAuth)
	router.CompanyRouter(r, companyService, jwtAuth)
	router.JobRouter(r, jobService, jwtAuth)
	router.TruckRouter(r, truckService, jwtAuth)

	log.Println("Starting server on :8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
