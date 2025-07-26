package router

import (
	"moveshare/internal/handlers/admin"

	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine, jwtAuth service.JWTAuth, adminService service.AdminService) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AdminMiddleware(jwtAuth, adminService))
	{
		adminGroup.GET("/users/count", admin.GetUserCount(adminService))
		adminGroup.GET("/conversations/count", admin.GetChatConversationCount(adminService))
		adminGroup.GET("/users", admin.GetUsersList(adminService))
		adminGroup.GET("/jobs", admin.GetAllJobs(adminService))
		adminGroup.PATCH("/user/:userID/status", admin.ChangeUserStatus(adminService))
	}
}
