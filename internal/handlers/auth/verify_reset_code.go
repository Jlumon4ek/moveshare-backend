package auth

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifyResetCode handles verification of password reset code
// @Summary Verify password reset code
// @Description Verify if the provided reset code is valid and not expired
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.VerifyResetCodeRequest true "Email and code"
// @Success 200 {object} models.VerifyResetCodeResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} models.VerifyResetCodeResponse
// @Router /verify-reset-code [post]
func VerifyResetCode(passwordResetService service.PasswordResetService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.VerifyResetCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		err := passwordResetService.VerifyResetCode(c.Request.Context(), req.Email, req.Code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.VerifyResetCodeResponse{
				Valid:   false,
				Message: "Invalid or expired reset code",
			})
			return
		}

		c.JSON(http.StatusOK, models.VerifyResetCodeResponse{
			Valid:   true,
			Message: "Reset code is valid",
		})
	}
}