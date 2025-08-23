package auth

import (
	"net/http"

	"moveshare/internal/service"
	"moveshare/internal/utils"

	"github.com/gin-gonic/gin"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} ChangePasswordResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Current password is incorrect"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /auth/change-password [post]
func ChangePassword(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		var request ChangePasswordRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate password length
		if len(request.NewPassword) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
			return
		}

		// Get user to verify current password
		user, err := userService.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
			return
		}

		if !userService.CheckPassword(request.CurrentPassword, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
			return
		}

		err = userService.UpdatePassword(userID, request.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, ChangePasswordResponse{
			Message: "Password changed successfully",
		})
	}
}
