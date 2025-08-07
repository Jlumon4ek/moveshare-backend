package company

import (
	"moveshare/internal/dto"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCompany godoc
// @Summary      Получить информацию о компании
// @Description  Возвращает компанию, связанную с текущим пользователем
// @Tags         Company
// @Security     BearerAuth
// @Produce      json
// @Router       /company/ [get]
// @Success      200  {object}  dto.CompanyResponse
func GetCompany(service service.CompanyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
			return
		}

		company, err := service.GetCompany(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get company",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, dto.NewCompanyResponse(company))
	}
}
