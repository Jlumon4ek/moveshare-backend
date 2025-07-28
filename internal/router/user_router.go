package router

import (
	"moveshare/internal/handlers/user"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(r gin.IRouter, userService service.UserService, jwtAuth service.JWTAuth) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/refresh-token", user.RefreshToken(userService, jwtAuth))
		authGroup.POST("/sign-up", user.SignUp(userService))
		authGroup.POST("/sign-in", user.SignIn(userService, jwtAuth))
	}
}
