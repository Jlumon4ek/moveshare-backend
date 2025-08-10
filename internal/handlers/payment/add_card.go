// internal/handlers/payment/add_card.go
package payment

import (
	"errors"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddCard godoc
// @Summary      Add a payment card
// @Description  Adds a new payment card for the authenticated user using Stripe Payment Method
// @Tags         Payment
// @Security     BearerAuth
// @Param        card body models.AddCardRequest true "Payment method data"
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.AddCardResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /payment/cards [post]
func AddCard(paymentService service.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Парсим тело запроса
		var req models.AddCardRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Валидируем payment method ID
		if err := validatePaymentMethodID(req.PaymentMethodID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Добавляем карту
		paymentMethod, err := paymentService.AddPaymentMethod(c.Request.Context(), userID, req.PaymentMethodID)
		if err != nil {
			// Проверяем тип ошибки для соответствующего HTTP статуса
			if strings.Contains(err.Error(), "already added") {
				c.JSON(http.StatusConflict, gin.H{
					"error":   "Payment method already added",
					"details": err.Error(),
				})
				return
			}

			if strings.Contains(err.Error(), "unsupported payment method type") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Unsupported payment method type",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to add payment method",
				"details": err.Error(),
			})
			return
		}

		// Формируем ответ
		response := models.AddCardResponse{
			PaymentMethod: models.PaymentMethodResponse{
				ID:           paymentMethod.ID,
				CardLast4:    paymentMethod.CardLast4,
				CardBrand:    paymentMethod.CardBrand,
				CardExpMonth: paymentMethod.CardExpMonth,
				CardExpYear:  paymentMethod.CardExpYear,
				IsDefault:    paymentMethod.IsDefault,
				CreatedAt:    paymentMethod.CreatedAt,
			},
			Message: "Payment method added successfully",
			Success: true,
		}

		c.JSON(http.StatusCreated, response)
	}
}

// validatePaymentMethodID валидирует Stripe Payment Method ID
func validatePaymentMethodID(paymentMethodID string) error {
	if paymentMethodID == "" {
		return errors.New("payment_method_id is required")
	}

	// Stripe Payment Method ID должен начинаться с "pm_"
	if !strings.HasPrefix(paymentMethodID, "pm_") {
		return errors.New("invalid payment method ID format. Must start with 'pm_'")
	}

	// Минимальная длина Stripe ID
	if len(paymentMethodID) < 20 {
		return errors.New("invalid payment method ID length")
	}

	return nil
}
