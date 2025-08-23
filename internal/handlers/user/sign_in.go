package user

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignIn handles user authentication
// @Summary Authenticate user
// @Description Authenticates user with email/username and password, returns JWT tokens
// @Tags Auth
// @Param request body models.SignInRequest true "User login data"
// @Router /auth/sign-in [post]
func SignIn(userService service.UserService, jwtAuth service.JWTAuth, sessionService service.SessionService) gin.HandlerFunc {
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

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid credentials"})
			return
		}

		// Create session record first
		userAgent := c.GetHeader("User-Agent")
		clientIP := utils.GetClientIP(c.Request)
		deviceInfo := utils.ParseUserAgent(userAgent)
		locationInfo := utils.GetLocationInfo(clientIP)

		sessionRequest := &models.CreateSessionRequest{
			UserAgent:    userAgent,
			IPAddress:    clientIP,
			DeviceInfo:   deviceInfo,
			LocationInfo: locationInfo,
		}

		// Create session with temporary tokens
		session, err := sessionService.CreateSession(c.Request.Context(), sessionRequest, user.ID, "temp", "temp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create session"})
			return
		}

		// Generate tokens with session ID
		accessToken, err := jwtAuth.GenerateAccessToken(user.ID, user.Username, user.Email, user.Role, session.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate access token"})
			return
		}

		refreshToken, err := jwtAuth.GenerateRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate refresh token"})
			return
		}

		// Update session with actual tokens
		err = sessionService.UpdateSessionTokens(c.Request.Context(), session.ID, accessToken, refreshToken)
		if err != nil {
			// Log error but don't fail the login
			// Session token update is not critical for login process
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
