package admin

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetJobsList handles getting a paginated list of jobs
// @Summary Get list of jobs
// @Description Gets a paginated list of jobs with limit, offset and optional status filter
// @Tags Admin
// @Produce json
// @Param limit query int false "Limit number of jobs returned" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Param status query string false "Filter by job status. Single status or comma-separated multiple (claimed, active, pending, canceled, completed)"
// @Success 200 {object} models.PaginatedJobsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/jobs [get]
// @Security     BearerAuth
func GetJobsList(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")
		statusParam := c.Query("status")
		
		var statuses []string
		// Parse and validate status parameter if provided
		if statusParam != "" {
			statusList := strings.Split(statusParam, ",")
			validStatuses := map[string]bool{
				"claimed":   true,
				"active":    true,
				"pending":   true,
				"canceled":  true,
				"completed": true,
			}
			
			for _, status := range statusList {
				status = strings.TrimSpace(status)
				if status == "" {
					continue
				}
				if !validStatuses[status] {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Invalid status parameter",
						"details": "Each status must be one of: claimed, active, pending, canceled, completed",
					})
					return
				}
				statuses = append(statuses, status)
			}
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid limit parameter",
				"details": err.Error(),
			})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid offset parameter",
				"details": err.Error(),
			})
			return
		}

		jobs, err := adminService.GetJobsList(c.Request.Context(), limit, offset, statuses)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get jobs list",
				"details": err.Error(),
			})
			return
		}

		total, err := adminService.GetJobsListTotal(c.Request.Context(), statuses)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get jobs total count",
				"details": err.Error(),
			})
			return
		}

		page := offset/limit + 1
		if offset == 0 && limit > 0 {
			page = 1
		}

		response := models.PaginatedJobsResponse{
			Jobs:   jobs,
			Total:  total,
			Page:   page,
			Limit:  limit,
			Offset: offset,
		}

		c.JSON(http.StatusOK, response)
	}
}