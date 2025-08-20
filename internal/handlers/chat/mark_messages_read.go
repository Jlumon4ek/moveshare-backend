package chat

import (
	"context"
	"log"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MarkMessagesAsRead godoc
// @Summary      Mark messages as read
// @Description  Mark all unread messages in a chat as read for the current user
// @Tags         Chat
// @Security     BearerAuth
// @Param        chatId path      int     true  "Chat ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chats/{chatId}/mark-read [post]
func MarkMessagesAsRead(chatService service.ChatService, notificationService service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		chatIDStr := c.Param("chatId")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil || chatID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
			return
		}

		// Verify user is participant of this chat
		isParticipant, err := chatService.IsUserParticipant(c.Request.Context(), chatID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to verify chat access",
				"details": err.Error(),
			})
			return
		}

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. You are not a participant of this chat"})
			return
		}

		// Mark messages as read
		err = chatService.MarkMessagesAsRead(c.Request.Context(), chatID, userID)
		if err != nil {
			log.Printf("Failed to mark messages as read for chat %d, user %d: %v", chatID, userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to mark messages as read",
				"details": err.Error(),
			})
			return
		}

		log.Printf("Successfully marked messages as read for chat %d, user %d", chatID, userID)
		
		// Отправляем обновление счетчика непрочитанных сообщений (асинхронно)
		if notificationService != nil {
			go func() {
				ctx := context.Background()
				unreadCount, err := chatService.GetUserUnreadCount(ctx, userID)
				if err != nil {
					log.Printf("Failed to get updated unread count for user %d: %v", userID, err)
				} else {
					notificationService.NotifyUnreadCountChange(userID, unreadCount)
				}
			}()
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Messages marked as read",
		})
	}
}