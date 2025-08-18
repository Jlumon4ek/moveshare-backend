package auth

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResetPassword handles password reset with code
// @Summary Reset password with code
// @Description Reset user password using email, code and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} models.ResetPasswordResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reset-password [post]
func ResetPassword(passwordResetService service.PasswordResetService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ResetPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		err := passwordResetService.ResetPassword(c.Request.Context(), req.Email, req.Code, req.NewPassword)
		if err != nil {
			if err.Error() == "invalid or expired reset code" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid or expired reset code",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to reset password",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ResetPasswordResponse{
			Message: "Password has been reset successfully",
		})
	}
}