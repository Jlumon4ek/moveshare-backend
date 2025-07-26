package company

import (
	"context"
	"errors"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetCompany godoc
// @Summary      Получить информацию о компании
// @Description  Возвращает компанию, связанную с текущим пользователем
// @Tags         Company
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
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusOK, &models.Company{})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get company",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, company)
	}
}
