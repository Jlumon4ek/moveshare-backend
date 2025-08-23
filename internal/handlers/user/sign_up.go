package user

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignUp handles user registration
// @Summary Register a new user
// @Description Creates a new user account with username, email, password and email verification code
// @Tags Auth
// @Param request body models.SignUpRequest true "User registration data"
// @Router /auth/sign-up [post]
func SignUp(userService service.UserService, emailVerificationService service.EmailVerificationService) gin.HandlerFunc {
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

		c.JSON(http.StatusCreated, models.SignUpResponse{Message: "User created successfully"})
	}
}
