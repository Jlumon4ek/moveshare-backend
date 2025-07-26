package admin

import (
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserCount handles getting user count
// @Summary Get total user count
// @Description Gets the total number of users in the system
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} map[string]string
// @Router /admin/users/count [get]
// @Security     BearerAuth
func GetUserCount(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCount, err := adminService.GetUserCount(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_count": userCount})
	}
}
