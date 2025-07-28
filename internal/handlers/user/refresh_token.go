package user

import (
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Validates refresh token and returns a new access token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body RefreshTokenRequest true "Refresh token payload"
// @Success      200 {object} map[string]interface{} "user_id and access_token"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      401 {object} map[string]string "Invalid or expired refresh token"
// @Failure      500 {object} map[string]string "Server error"
// @Router       /auth/refresh-token [post]
func RefreshToken(userService service.UserService, jwtAuth service.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
			return
		}

		userID, err := jwtAuth.ValidateToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		userInfo, err := userService.FindUserByID(c.Request.Context(), userID)
		if err != nil || userInfo == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   error.Error(err),
				"details": "User not found or invalid token",
			})
			return
		}

		accessToken, err := jwtAuth.GenerateAccessToken(userID, userInfo.Username, userInfo.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":      userID,
			"access_token": accessToken,
		})
	}
}
