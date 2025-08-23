package user

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMyProfile godoc
// @Summary      Get current user profile
// @Description  Returns the profile information of the currently authenticated user including updated_at
// @Tags         User
// @Produce      json
// @Success      200 {object} models.User "user profile"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Server error"
// @Router       /user/my-profile [get]
// @Security     BearerAuth
func GetMyProfile(userService service.UserService) gin.HandlerFunc {
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

		// Remove sensitive information before sending to client
		userInfo.Password = ""

		c.JSON(http.StatusOK, userInfo)
	}
}