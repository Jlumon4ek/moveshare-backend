package router

import (
	"moveshare/internal/handlers/chat"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func ChatRouter(r gin.IRouter, chatService service.ChatService, jwtAuth service.JWTAuth, hub *chat.Hub) {
	chatGroup := r.Group("/chats")
	chatGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		chatGroup.GET("/", chat.GetUserChats(chatService))
		chatGroup.GET("/:chatId/messages", chat.GetChatMessages(chatService))
		chatGroup.POST("/:chatId/messages", chat.SendMessage(chatService, hub))
	}

	r.GET("/chats/:chatId/ws", chat.WebSocketChat(hub, jwtAuth))
}
