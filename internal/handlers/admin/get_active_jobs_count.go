package admin

import (
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetActiveJobsCount handles getting active jobs count
// @Summary Get active jobs count
// @Description Gets the total number of active jobs (not completed) in the system
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} map[string]string
// @Router /admin/jobs/active/count [get]
// @Security     BearerAuth
func GetActiveJobsCount(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		activeJobsCount, err := adminService.GetActiveJobsCount(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active jobs count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"active_jobs_count": activeJobsCount})
	}
}