package user

import (
	"moveshare/internal/schemas"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignIn handles user authentication
// @Summary Authenticate user
// @Description Authenticates user with email/username and password, returns JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body schemas.SignInRequest true "User login data"
// @Success 200 {object} schemas.SignInResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /auth/sign-in [post]
func SignIn(userService service.UserService, jwtAuth service.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req schemas.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "Invalid request body"})
			return
		}

		if req.Identifier == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "Identifier and password are required"})
			return
		}

		user, err := userService.FindUserByEmailOrUsername(c.Request.Context(), req.Identifier)
		if err != nil {
			c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{Error: "Invalid credentials"})
			return
		}

		accessToken, err := jwtAuth.GenerateAccessToken(user.ID, user.Username, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: "Failed to generate access token"})
			return
		}

		refreshToken, err := jwtAuth.GenerateRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, schemas.SignInResponse{
			UserID:       user.ID,
			Username:     user.Username,
			Email:        user.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
