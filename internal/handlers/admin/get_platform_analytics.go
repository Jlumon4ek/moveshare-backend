package admin

import (
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPlatformAnalytics handles getting platform analytics data
// @Summary Get platform analytics
// @Description Gets platform analytics including top companies and busiest routes
// @Tags Admin
// @Produce json
// @Param days query int false "Number of days to look back" default(30)
// @Success 200 {object} models.PlatformAnalytics
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/analytics [get]
// @Security     BearerAuth
func GetPlatformAnalytics(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get days parameter from query, default to 30
		daysStr := c.Query("days")
		days := 30
		if daysStr != "" {
			parsedDays, err := strconv.Atoi(daysStr)
			if err != nil || parsedDays < 1 || parsedDays > 365 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
				return
			}
			days = parsedDays
		}

		analytics, err := adminService.GetPlatformAnalytics(c.Request.Context(), days)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platform analytics"})
			return
		}

		c.JSON(http.StatusOK, analytics)
	}
}