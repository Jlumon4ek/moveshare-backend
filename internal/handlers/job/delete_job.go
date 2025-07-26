package job

import (
	"fmt"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteJob godoc
// @Summary      Delete a job
// @Description  User deletes a job with given ID
// @Tags         Jobs
// @Param        jobID path int true "Job ID"
// @Security     BearerAuth
// @Router       /jobs/{jobID} [delete]
func DeleteJob(jobService service.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		jobIDParam := c.Param("jobID")
		if jobIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing jobID parameter"})
			return
		}
		var jobID int64
		if _, err := fmt.Sscan(jobIDParam, &jobID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jobID parameter"})
			return
		}

		if err := jobService.DeleteJob(c.Request.Context(), userID, jobID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete job: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
	}
}
