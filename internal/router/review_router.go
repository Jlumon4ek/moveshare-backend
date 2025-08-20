package router

import (
	"moveshare/internal/handlers/review"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupReviewRoutes(router gin.IRouter, reviewHandler *review.ReviewHandler, jwtService service.JWTAuth) {
	reviewRoutes := router.Group("/reviews")
	reviewRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		// POST /reviews - Создать отзыв
		reviewRoutes.POST("/", reviewHandler.CreateReview)

		// GET /reviews/user/:id - Получить отзывы пользователя
		reviewRoutes.GET("/user/:id", reviewHandler.GetUserReviews)

		// GET /reviews/stats/:id - Получить статистику рейтинга пользователя
		reviewRoutes.GET("/stats/:id", reviewHandler.GetUserRatingStats)
		
		// GET /reviews/average/:id - Получить среднюю оценку пользователя
		reviewRoutes.GET("/average/:id", reviewHandler.GetUserAverageRating)
		
		// GET /reviews/job/:id/check - Проверить существование отзыва для работы
		reviewRoutes.GET("/job/:id/check", reviewHandler.CheckJobReviewExists)
	}
}
