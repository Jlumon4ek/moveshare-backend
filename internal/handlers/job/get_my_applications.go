package job

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMyApplications godoc
// @Summary      Get my job applications
// @Description  Returns a list of job applications submitted by the authenticated user
// @Tags         Jobs
// @Router       /jobs/applications/my [get]
// @Security     BearerAuth
func GetMyApplications(jobService service.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		applications, err := jobService.GetMyApplications(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
			return
		}

		c.JSON(http.StatusOK, applications)
	}
}
