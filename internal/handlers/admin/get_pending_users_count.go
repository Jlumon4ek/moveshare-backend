package admin

import (
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPendingUsersCount handles getting pending users count
// @Summary Get pending users count
// @Description Gets the total number of users with "On Waiting" status
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} map[string]string
// @Router /admin/users/pending/count [get]
// @Security     BearerAuth
func GetPendingUsersCount(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		pendingUsersCount, err := adminService.GetPendingUsersCount(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending users count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"pending_users_count": pendingUsersCount})
	}
}