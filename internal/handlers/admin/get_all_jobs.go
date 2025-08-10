package admin

// import (
// 	"moveshare/internal/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // GetActiveJobs handles getting a paginated list of active jobs
// // @Summary Get list of active jobs
// // @Description Gets a paginated list of active jobs with limit and offset
// // @Tags Admin
// // @Produce json
// // @Param limit query int false "Limit number of jobs returned" default(10)
// // @Param offset query int false "Offset for pagination" default(0)
// // @Success 200 {object} []models.Job
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /admin/jobs [get]
// // @Security     BearerAuth
// func GetAllJobs(adminService service.AdminService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		limitStr := c.DefaultQuery("limit", "10")
// 		offsetStr := c.DefaultQuery("offset", "0")

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

// 		jobs, err := adminService.GetActiveJobs(c.Request.Context(), limit, offset)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active jobs"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"jobs": jobs})
// 	}
// }
