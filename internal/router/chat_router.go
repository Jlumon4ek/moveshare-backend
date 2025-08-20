package router

import (
	"moveshare/internal/handlers/chat"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(r gin.IRouter, chatService service.ChatService, jobService service.JobService, jwtAuth service.JWTAuth, hub *chat.Hub, notificationService service.NotificationService) {
	chatGroup := r.Group("/chats")
	chatGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		chatGroup.GET("/", chat.GetUserChats(chatService))
		chatGroup.GET("/:chatId/messages", chat.GetChatMessages(chatService))
		chatGroup.POST("/:chatId/messages", chat.SendMessage(chatService, hub, notificationService))
		chatGroup.POST("/:chatId/mark-read", chat.MarkMessagesAsRead(chatService, notificationService))
		chatGroup.POST("/", chat.CreateChat(chatService, jobService)) // ✅ Создание чата
	}

	r.GET("/chats/:chatId/ws", chat.WebSocketChat(hub, jwtAuth))
}
