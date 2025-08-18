package auth

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ForgotPassword handles password reset request
// @Summary Request password reset
// @Description Send password reset code to user's email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Email address"
// @Success 200 {object} models.ForgotPasswordResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /forgot-password [post]
func ForgotPassword(passwordResetService service.PasswordResetService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ForgotPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		err := passwordResetService.SendResetCode(c.Request.Context(), req.Email)
		if err != nil {
			if err.Error() == "user not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User with this email does not exist",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to send reset code",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ForgotPasswordResponse{
			Message: "Password reset code has been sent to your email",
		})
	}
}