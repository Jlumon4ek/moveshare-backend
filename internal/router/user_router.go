package router

import (
	"moveshare/internal/handlers/auth"
	"moveshare/internal/handlers/session"
	"moveshare/internal/handlers/user"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(r gin.IRouter, userService service.UserService, minioService *service.MinioService, jwtAuth service.JWTAuth, passwordResetService service.PasswordResetService, sessionService service.SessionService) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/refresh-token", user.RefreshToken(userService, jwtAuth))
		authGroup.POST("/sign-up", user.SignUp(userService))
		authGroup.POST("/sign-in", user.SignIn(userService, jwtAuth, sessionService))
		// Logout requires authentication
		authGroup.Use(middleware.AuthMiddleware(jwtAuth))
		authGroup.Use(middleware.SessionValidationMiddleware(sessionService))
		authGroup.POST("/logout", auth.Logout(sessionService))
	}

	// Password reset routes (public)
	r.POST("/forgot-password", auth.ForgotPassword(passwordResetService))
	r.POST("/verify-reset-code", auth.VerifyResetCode(passwordResetService))
	r.POST("/reset-password", auth.ResetPassword(passwordResetService))
	profilePhotoHandler := user.NewProfilePhotoHandler(userService, minioService)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware(jwtAuth))
	userGroup.Use(middleware.SessionValidationMiddleware(sessionService))
	{
		userGroup.GET("/my-status", user.GetMyStatus(userService))
		userGroup.GET("/my-profile", user.GetMyProfile(userService))
		userGroup.POST("/upload-profile-photo", profilePhotoHandler.UploadProfilePhoto)
		userGroup.GET("/profile-photo/:user_id", profilePhotoHandler.GetProfilePhoto)
		userGroup.DELETE("/profile-photo", profilePhotoHandler.DeleteProfilePhoto)
		userGroup.POST("/change-password", auth.ChangePassword(userService))

		// Session management routes
		userGroup.GET("/active-sessions", session.GetActiveSessions(sessionService))
		userGroup.DELETE("/sessions/:session_id/terminate", session.TerminateSession(sessionService))
		userGroup.DELETE("/sessions/terminate-all", session.TerminateAllSessions(sessionService))
	}

}
