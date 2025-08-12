// internal/router/payment_router.go
package router

import (
	"moveshare/internal/handlers/payment"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func PaymentRouter(r gin.IRouter, paymentService service.PaymentService, jwtAuth service.JWTAuth) {
	paymentGroup := r.Group("/payment")
	paymentGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		// Payment Methods (Cards)
		paymentGroup.POST("/cards", payment.AddCard(paymentService))
		paymentGroup.GET("/cards", payment.GetUserCards(paymentService))
		paymentGroup.DELETE("/cards/:cardId", payment.DeleteCard(paymentService))
		paymentGroup.PATCH("/cards/:cardId/default", payment.SetDefaultCard(paymentService))

		paymentGroup.POST("/create-intent", payment.CreatePayment(paymentService))
		paymentGroup.POST("/confirm-payment", payment.ConfirmPayment(paymentService))
		paymentGroup.GET("/history", payment.GetPaymentHistory(paymentService))
	}

	// Webhook endpoint (без аутентификации)
	// r.POST("/payment/webhook", payment.StripeWebhook(paymentService))
}
