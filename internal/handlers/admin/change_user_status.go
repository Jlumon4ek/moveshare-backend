package admin

import (
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ChangeUserStatus godoc
// @Summary      Change user status
// @Description  Changes the status of a user by their ID
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        userID   path      int    true  "User ID"
// @Param        status   query     string true  "New status for the user"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /admin/user/{userID}/status [patch]
// @Security     BearerAuth
func ChangeUserStatus(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("userID")
		newStatus := c.Query("status")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if newStatus == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New status must be provided"})
			return
		}

		err = adminService.ChangeUserStatus(c.Request.Context(), userID, newStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change user status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
	}
}
