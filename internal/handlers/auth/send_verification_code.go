package auth

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendVerificationCode godoc
// @Summary Send email verification code
// @Description Sends a 6-digit verification code to the specified email address
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.SendVerificationCodeRequest true "Email address"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/send-verification-code [post]
func SendVerificationCode(emailVerificationService service.EmailVerificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.SendVerificationCodeRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		err := emailVerificationService.SendVerificationCode(c.Request.Context(), request.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Verification code sent successfully",
		})
	}
}