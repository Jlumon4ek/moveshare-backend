package admin

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSystemSettings handles getting system settings
// @Summary Get system settings
// @Description Gets the current system settings
// @Tags Admin
// @Produce json
// @Success 200 {object} models.SystemSettings
// @Failure 500 {object} map[string]string
// @Router /admin/settings [get]
// @Security     BearerAuth
func GetSystemSettings(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		settings, err := adminService.GetSystemSettings(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get system settings"})
			return
		}

		c.JSON(http.StatusOK, settings)
	}
}

// UpdateSystemSettings handles updating system settings
// @Summary Update system settings
// @Description Updates the system settings
// @Tags Admin
// @Accept json
// @Produce json
// @Param settings body models.SystemSettings true "System settings"
// @Success 200 {object} models.SystemSettings
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/settings [put]
// @Security     BearerAuth
func UpdateSystemSettings(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var settings models.SystemSettings
		if err := c.ShouldBindJSON(&settings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate values
		if settings.CommissionRate < 0 || settings.CommissionRate > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Commission rate must be between 0 and 100"})
			return
		}

		if settings.NewUserApproval != "manual" && settings.NewUserApproval != "auto" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New user approval must be 'manual' or 'auto'"})
			return
		}

		if settings.MinimumPayout < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum payout must be positive"})
			return
		}

		if settings.JobExpirationDays < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Job expiration days must be at least 1"})
			return
		}

		err := adminService.UpdateSystemSettings(c.Request.Context(), &settings)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update system settings"})
			return
		}

		c.JSON(http.StatusOK, settings)
	}
}