package company

import (
	"context"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCompany godoc
// @Summary      Получить информацию о компании
// @Description  Возвращает компанию, связанную с текущим пользователем
// @Tags         company
// @Security     BearerAuth
// @Produce      json
// @Router       /company/ [get]
func GetCompany(service service.CompanyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		company, err := service.GetCompany(context.Background(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get company",
				"details": err.Error(),
			})
			return
		}
		if company == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}
		c.JSON(http.StatusOK, company)
	}
}
