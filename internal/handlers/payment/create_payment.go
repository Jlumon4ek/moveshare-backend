// internal/handlers/payment/create_payment.go
package payment

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreatePayment godoc
// @Summary      Create a payment
// @Description  Creates a payment intent for a job using user's saved card
// @Tags         Payment
// @Security     BearerAuth
// @Param        payment body models.CreatePaymentRequest true "Payment data"
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.CreatePaymentResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /payment/create-intent [post]
func CreatePayment(paymentService service.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Парсим тело запроса
		var req models.CreatePaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Создаем платеж
		response, err := paymentService.CreatePayment(c.Request.Context(), userID, &req)
		if err != nil {
			// Проверяем тип ошибки для соответствующего HTTP статуса
			if strings.Contains(err.Error(), "payment method not found") ||
				strings.Contains(err.Error(), "no default payment method found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "Payment method not found",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create payment",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}
