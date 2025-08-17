package user

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignIn handles user authentication
// @Summary Authenticate user
// @Description Authenticates user with email/username and password, returns JWT tokens
// @Tags Auth
// @Param request body models.SignInRequest true "User login data"
// @Router /auth/sign-in [post]
func SignIn(userService service.UserService, jwtAuth service.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
			return
		}

		if req.Identifier == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Identifier and password are required"})
			return
		}

		user, err := userService.FindUserByEmailOrUsername(c.Request.Context(), req.Identifier)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid credentials"})
			return
		}

		accessToken, err := jwtAuth.GenerateAccessToken(user.ID, user.Username, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate access token"})
			return
		}

		refreshToken, err := jwtAuth.GenerateRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, models.SignInResponse{
			UserID:       user.ID,
			Username:     user.Username,
			Email:        user.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
