package payment

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteCard godoc
// @Summary      Delete a payment card
// @Description  Deletes a payment card for the authenticated user
// @Tags         Payment
// @Security     BearerAuth
// @Param        cardId path int true "Card ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /payment/cards/{cardId} [delete]
func DeleteCard(paymentService service.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Парсим ID карты из URL
		cardIDStr := c.Param("cardId")
		cardID, err := strconv.ParseInt(cardIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid card ID",
				"details": "Card ID must be a valid number",
			})
			return
		}

		// Удаляем карту
		err = paymentService.DeletePaymentMethod(c.Request.Context(), userID, cardID)
		if err != nil {
			// Проверяем тип ошибки для соответствующего HTTP статуса
			if err.Error() == "sql: no rows in result set" ||
				err.Error() == "payment method not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "Payment method not found",
					"details": "The specified card does not exist or does not belong to you",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to delete payment method",
				"details": err.Error(),
			})
			return
		}

		// Успешный ответ
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment method deleted successfully",
			"success": true,
		})
	}
}
