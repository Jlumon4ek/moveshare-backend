// internal/handlers/payment/get_cards.go
package payment

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserCards godoc
// @Summary      Get user payment cards
// @Description  Gets all active payment cards for the authenticated user
// @Tags         Payment
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /payment/cards [get]
func GetUserCards(paymentService service.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Получаем карты пользователя
		paymentMethods, err := paymentService.GetUserPaymentMethods(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get payment methods",
				"details": err.Error(),
			})
			return
		}

		// Формируем ответ
		response := gin.H{
			"cards":       paymentMethods,
			"total_cards": len(paymentMethods),
			"success":     true,
		}

		c.JSON(http.StatusOK, response)
	}
}
