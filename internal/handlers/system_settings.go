package handlers

import (
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCommissionRate handles getting commission rate for users
// @Summary Get platform commission rate
// @Description Gets the current platform commission rate that users can see
// @Tags System
// @Produce json
// @Success 200 {object} map[string]float64
// @Failure 500 {object} map[string]string
// @Router /commission-rate [get]
func GetCommissionRate(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		settings, err := adminService.GetSystemSettings(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get commission rate"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"commission_rate": settings.CommissionRate,
		})
	}
}