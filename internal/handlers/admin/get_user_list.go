package admin

import (
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUsersList handles getting a paginated list of users
// @Summary Get list of users
// @Description Gets a paginated list of users with limit and offset
// @Tags Admin
// @Produce json
// @Param limit query int false "Limit number of users returned" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} []models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func GetUsersList(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

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

		users, err := adminService.GetUsersList(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get users list",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
