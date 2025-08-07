package company

import (
	"context"
	"moveshare/internal/dto"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PatchCompany godoc
// @Summary      Обновить информацию о компании
// @Description  Частично обновляет компанию, связанную с текущим пользователем
// @Tags         Company
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        company  body      dto.UpdateCompanyRequest  true  "Данные для обновления компании"
// @Router       /company/ [put]
func PatchCompany(service service.CompanyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		var req dto.UpdateCompanyRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		if err := service.UpdateCompany(context.Background(), userID, req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company", "details": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
