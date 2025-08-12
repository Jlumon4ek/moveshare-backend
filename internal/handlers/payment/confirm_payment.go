// internal/handlers/payment/confirm_payment.go
package payment

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ConfirmPayment godoc
// @Summary      Confirm a payment
// @Description  Confirms a payment intent after user completes 3D Secure or other authentication
// @Tags         Payment
// @Security     BearerAuth
// @Param        confirmation body models.ConfirmPaymentRequest true "Payment confirmation data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ConfirmPaymentResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /payment/confirm-payment [post]
func ConfirmPayment(paymentService service.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Парсим тело запроса
		var req models.ConfirmPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Подтверждаем платеж
		response, err := paymentService.ConfirmPayment(c.Request.Context(), req.PaymentIntentID)
		if err != nil {
			if err.Error() == "payment not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "Payment not found",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to confirm payment",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
