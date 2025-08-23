package auth

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifyEmailCode godoc
// @Summary Verify email verification code
// @Description Verifies the 6-digit email verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.VerifyEmailCodeRequest true "Email and verification code"
// @Success 200 {object} models.VerifyEmailCodeResponse "Code verified successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid or expired code"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/verify-email-code [post]
func VerifyEmailCode(emailVerificationService service.EmailVerificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.VerifyEmailCodeRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		err := emailVerificationService.VerifyEmailCode(c.Request.Context(), request.Email, request.Code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification code"})
			return
		}

		c.JSON(http.StatusOK, models.VerifyEmailCodeResponse{
			Message: "Email verified successfully",
			Valid:   true,
		})
	}
}