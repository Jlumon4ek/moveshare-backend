package job

import (
	"moveshare/internal/dto"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMyJobs godoc
// @Summary      Get jobs created by current user
// @Description  Returns a list of jobs posted by the authenticated user
// @Tags         Jobs
// @Security     BearerAuth
// @Router       /jobs/my [get]
func GetMyJobs(jobService service.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		jobModels, err := jobService.GetUserJobs(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get user jobs",
				"details": err.Error(),
			})
			return
		}

		var jobs []dto.JobResponse
		for _, j := range jobModels {
			jobs = append(jobs, dto.NewJobResponse(j))
		}

		c.JSON(http.StatusOK, gin.H{"jobs": jobs})
	}
}
