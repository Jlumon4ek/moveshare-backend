package router

import (
	"moveshare/internal/handlers/verification"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func VerificationRouter(r gin.IRouter, verificationService service.VerificationService, jwtAuth service.JWTAuth) {
	verificationGroup := r.Group("/verification")
	verificationGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		verificationGroup.POST("/", verification.UploadVerificationFile(verificationService))
		verificationGroup.GET("/files", verification.GetVerificationFile(verificationService))
	}
}
