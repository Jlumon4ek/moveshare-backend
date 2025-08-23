package user

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignUp handles user registration
// @Summary Register a new user
// @Description Creates a new user account with username, email, password and email verification code
// @Tags Auth
// @Param request body models.SignUpRequest true "User registration data"
// @Router /auth/sign-up [post]
func SignUp(userService service.UserService, emailVerificationService service.EmailVerificationService, jwtAuth service.JWTAuth, sessionService service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignUpRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
			return
		}

		if req.Username == "" || req.Email == "" || req.Password == "" || req.VerificationCode == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Username, email, password and verification code are required"})
			return
		}

		// Verify email verification code first
		err := emailVerificationService.VerifyEmailCode(c.Request.Context(), req.Email, req.VerificationCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid or expired verification code"})
			return
		}

		existingUser, err := userService.FindUserByEmailOrUsername(c.Request.Context(), req.Email)
		if err == nil && existingUser != nil {
			c.JSON(http.StatusConflict, models.ErrorResponse{Error: "User with this email already exists"})
			return
		}

		existingUser, err = userService.FindUserByEmailOrUsername(c.Request.Context(), req.Username)
		if err == nil && existingUser != nil {
			c.JSON(http.StatusConflict, models.ErrorResponse{Error: "User with this username already exists"})
			return
		}

		user := &models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := userService.CreateUser(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create user"})
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
			// Log error but don't fail the registration process
		}

		c.JSON(http.StatusCreated, models.SignUpResponse{
			UserID:       user.ID,
			Username:     user.Username,
			Email:        user.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
