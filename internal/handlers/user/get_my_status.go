package user

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMyStatus godoc
// @Summary      Get current user status
// @Description  Returns the status of the currently authenticated user
// @Tags         User
// @Produce      json
// @Success      200 {object} map[string]string "user status"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Server error"
// @Router       /user/my-status [get]
// @Security     BearerAuth
func GetMyStatus(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userInfo, err := userService.FindUserByID(c.Request.Context(), userID)
		if err != nil || userInfo == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": userInfo.Status,
		})
	}
}