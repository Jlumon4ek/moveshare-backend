package router

import (
	"moveshare/internal/handlers/auth"
	"moveshare/internal/handlers/user"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(r gin.IRouter, userService service.UserService, minioService *service.MinioService, jwtAuth service.JWTAuth, passwordResetService service.PasswordResetService) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/refresh-token", user.RefreshToken(userService, jwtAuth))
		authGroup.POST("/sign-up", user.SignUp(userService))
		authGroup.POST("/sign-in", user.SignIn(userService, jwtAuth))
	}

	// Password reset routes (public)
	r.POST("/forgot-password", auth.ForgotPassword(passwordResetService))
	r.POST("/verify-reset-code", auth.VerifyResetCode(passwordResetService))
	r.POST("/reset-password", auth.ResetPassword(passwordResetService))
	profilePhotoHandler := user.NewProfilePhotoHandler(userService, minioService)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		userGroup.GET("/my-status", user.GetMyStatus(userService))
		userGroup.POST("/upload-profile-photo", profilePhotoHandler.UploadProfilePhoto)
		userGroup.GET("/profile-photo/:user_id", profilePhotoHandler.GetProfilePhoto)
		userGroup.DELETE("/profile-photo", profilePhotoHandler.DeleteProfilePhoto)
	}

}
