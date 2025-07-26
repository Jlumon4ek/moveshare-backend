package job

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostNewJob godoc
// @Summary      Create new job
// @Description  Creates a new job for the authenticated user
// @Tags         Jobs
// @Param        job  body      models.PostNewJobRequest  true  "Job data"
// @Router       /jobs/post-new-job/ [post]
// @Security     BearerAuth
func PostNewJob(jobService service.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var job models.Job
		if err := c.ShouldBindJSON(&job); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		job.UserID = userID

		if err := jobService.CreateJob(c.Request.Context(), &job, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Job created successfully"})
	}
}
