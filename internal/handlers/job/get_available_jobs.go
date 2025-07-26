package job

// import (
// 	"moveshare/internal/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// func GetAvailableJobs(jobService service.JobService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		limitStr := c.DefaultQuery("limit", "10")
// 		offsetStr := c.DefaultQuery("offset", "0")

// 		filters := make(map[string]string)
// 		values := c.Request.URL.Query()
// 		if val := values.Get("pickup_location"); val != "" {
// 			filters["pickup_location"] = val
// 		}
// 		if val := values.Get("delivery_location"); val != "" {
// 			filters["delivery_location"] = val
// 		}
// 		if val := values.Get("pickup_date_start"); val != "" {
// 			filters["pickup_date_start"] = val
// 		}
// 		if val := values.Get("pickup_date_end"); val != "" {
// 			filters["pickup_date_end"] = val
// 		}
// 		if val := values.Get("truck_size"); val != "" {
// 			filters["truck_size"] = val
// 		}
// 		limit, err := strconv.Atoi(limitStr)
// 		if err != nil || limit < 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
// 			return
// 		}

// 		offset, err := strconv.Atoi(offsetStr)
// 		if err != nil || offset < 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
// 			return
// 		}

// 		jobs, err := jobService.GetAvailableJobs(c.Request.Context(), limit, offset)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get available jobs"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"jobs": jobs})
// 	}
// }
