package job

import (
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"moveshare/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetAvailableJobs godoc
// @Summary      Get available jobs
// @Description  Returns a list of available jobs based on filters and pagination
// @Tags         Jobs
// @Security     BearerAuth
// @Param        limit              query     int     false  "Limit (default 10)"
// @Param        offset             query     int     false  "Offset (default 0)"
// @Param        pickup_location    query     string  false  "Pickup location"
// @Param        delivery_location  query     string  false  "Delivery location"
// @Param        pickup_date_start  query     string  false  "Pickup date start (YYYY-MM-DD)"
// @Param        pickup_date_end    query     string  false  "Pickup date end (YYYY-MM-DD)"
// @Param        truck_size         query     string  false  "Truck size (e.g., small, medium, large)"
// @Router       /jobs/available/ [get]
func GetAvailableJobs(jobService service.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

		filters := make(map[string]string)
		values := c.Request.URL.Query()
		if val := values.Get("pickup_location"); val != "" {
			filters["pickup_location"] = val
		}
		if val := values.Get("delivery_location"); val != "" {
			filters["delivery_location"] = val
		}
		if val := values.Get("pickup_date_start"); val != "" {
			filters["pickup_date_start"] = val
		}
		if val := values.Get("pickup_date_end"); val != "" {
			filters["pickup_date_end"] = val
		}
		if val := values.Get("truck_size"); val != "" {
			filters["truck_size"] = val
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}

		jobs, err := jobService.GetAvailableJobs(c.Request.Context(), userID, filters, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get available jobs",
				"details": err.Error(),
			})
			return
		}

		if len(jobs) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"jobs":    []string{},
				"message": "No available jobs found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"jobs": jobs})
	}
}
