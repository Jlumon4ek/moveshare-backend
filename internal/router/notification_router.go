package router

import (
	"moveshare/internal/handlers/notifications"
	"moveshare/internal/middleware"
	"moveshare/internal/service"
	"moveshare/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(r gin.IRouter, jwtAuth service.JWTAuth, notificationHub *websocket.NotificationHub, notificationService service.NotificationService) {
	// WebSocket endpoint (no auth middleware needed as it handles JWT via query param)
	r.GET("/notifications/ws", notifications.WebSocketNotifications(notificationHub, jwtAuth))
	
	// REST API endpoints
	notificationGroup := r.Group("/notifications")
	notificationGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		// Get notifications
		notificationGroup.GET("/", notifications.GetNotifications(notificationService))
		notificationGroup.GET("/stats", notifications.GetNotificationStats(notificationService))
		
		// Mark as read
		notificationGroup.POST("/:id/read", notifications.MarkNotificationAsRead(notificationService))
		notificationGroup.POST("/read-all", notifications.MarkAllNotificationsAsRead(notificationService))
		
		// Delete notifications
		notificationGroup.DELETE("/:id", notifications.DeleteNotification(notificationService))
		notificationGroup.DELETE("/", notifications.DeleteAllNotifications(notificationService))
	}
}