// internal/handlers/payment/get_payment_history.go
package payment

import (
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

// GetPaymentHistory godoc
// @Summary      Get payment history
// @Description  Gets paginated payment history for the authenticated user
// @Tags         Payment
// @Security     BearerAuth
// @Param        limit query int false "Limit number of payments returned" default(10)
// @Param        offset query int false "Offset for pagination" default(0)
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /payment/history [get]
func GetPaymentHistory(paymentService service.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Парсим параметры пагинации
		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid limit parameter",
				"details": "Limit must be between 1 and 100",
			})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid offset parameter",
				"details": "Offset must be 0 or greater",
			})
			return
		}

		// Получаем историю платежей
		payments, err := paymentService.GetUserPayments(c.Request.Context(), userID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get payment history",
				"details": err.Error(),
			})
			return
		}

		// Формируем ответ
		response := gin.H{
			"payments": payments,
			"pagination": gin.H{
				"limit":  limit,
				"offset": offset,
				"count":  len(payments),
			},
			"success": true,
		}

		c.JSON(http.StatusOK, response)
	}
}
