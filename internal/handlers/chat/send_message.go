// internal/handlers/chat/send_message.go
package chat

import (
	"context"
	"errors"
	"log"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SendMessageRequest struct {
	MessageText string `json:"message_text" binding:"required" example:"Привет! Интересует ваше предложение."`
	MessageType string `json:"message_type,omitempty" example:"text"`
}

type SendMessageResponse struct {
	Message models.ChatMessageResponse `json:"message"`
	Success bool                       `json:"success"`
}

// SendMessage godoc
// @Summary      Send message to chat
// @Description  Sends a new message to the specified chat
// @Tags         Chat
// @Security     BearerAuth
// @Param        chatId path      int                 true  "Chat ID"
// @Param        message body     SendMessageRequest  true  "Message data"
// @Accept       json
// @Produce      json
// @Success      201  {object}  SendMessageResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chats/{chatId}/messages [post]
func SendMessage(chatService service.ChatService, hub *Hub, notificationService service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Получаем chatId из URL параметра
		chatIDStr := c.Param("chatId")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil || chatID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
			return
		}

		// Парсим тело запроса
		var req SendMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Валидируем сообщение
		if err := validateMessage(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверяем, что пользователь является участником чата
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

		// Проверяем активность чата
		isActive, err := chatService.IsChatActive(c.Request.Context(), chatID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check chat status",
				"details": err.Error(),
			})
			return
		}

		if !isActive {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send messages to inactive chat"})
			return
		}

		// Создаем объект сообщения
		message := &models.ChatMessage{
			ConversationID: chatID,
			SenderID:       userID,
			MessageText:    strings.TrimSpace(req.MessageText),
			MessageType:    getMessageType(req.MessageType),
		}

		// Отправляем сообщение
		createdMessage, err := chatService.SendMessage(c.Request.Context(), message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to send message",
				"details": err.Error(),
			})
			return
		}

		// ✅ ИСПРАВЛЕНИЕ: Отправляем через WebSocket всем ДРУГИМ участникам чата
		if hub != nil {
			// Создаем копию сообщения для WebSocket (для других пользователей)
			wsMessage := *createdMessage
			wsMessage.IsFromMe = false // Для получателей это НЕ их сообщение

			hub.BroadcastMessage(chatID, userID, &wsMessage)
		}

		// Отправляем push-уведомления участникам чата (асинхронно)
		if notificationService != nil {
			go func() {
				// Создаем новый контекст для горутины
				ctx := context.Background()
				log.Printf("Starting notification process for chat %d", chatID)
				participants, err := chatService.GetChatParticipants(ctx, chatID)
				if err != nil {
					log.Printf("Failed to get chat participants for chat %d: %v", chatID, err)
					return
				}

				log.Printf("Found %d participants in chat %d", len(participants), chatID)
				for _, participant := range participants {
					log.Printf("Processing participant %d (name: %s, role: %s)", participant.UserID, participant.UserName, participant.Role)
					
					// Не отправляем уведомление отправителю
					if participant.UserID == userID {
						log.Printf("Skipping sender %d", participant.UserID)
						continue
					}

					// Получаем обновленное количество непрочитанных сообщений для этого пользователя
					unreadCount, err := chatService.GetUserUnreadCount(ctx, participant.UserID)
					if err != nil {
						log.Printf("Failed to get unread count for user %d: %v", participant.UserID, err)
					} else {
						// Отправляем обновление счетчика через WebSocket
						notificationService.NotifyUnreadCountChange(participant.UserID, unreadCount)
					}

					// Проверяем, подключен ли пользователь к WebSocket чата
					if hub != nil && hub.IsUserConnectedToChat(chatID, participant.UserID) {
						log.Printf("User %d is connected to chat %d WebSocket, skipping push notification", participant.UserID, chatID)
						continue
					}

					log.Printf("Sending push notification to user %d", participant.UserID)
					// Отправляем уведомление пользователю
					notificationService.NotifyNewMessage(
						participant.UserID,
						chatID,
						createdMessage.SenderName,
						createdMessage.MessageText,
					)
				}
			}()
		}

		// ✅ ИСПРАВЛЕНИЕ: В REST ответе устанавливаем правильное значение для отправителя
		createdMessage.IsFromMe = true // Для отправителя это ЕГО сообщение

		// Обновляем время последней активности чата (асинхронно)
		go func() {
			_ = chatService.UpdateChatActivity(c.Request.Context(), chatID)
		}()

		c.JSON(http.StatusCreated, SendMessageResponse{
			Message: *createdMessage,
			Success: true,
		})
	}
}

// validateMessage валидирует сообщение
func validateMessage(req *SendMessageRequest) error {
	// Проверяем длину сообщения
	if len(strings.TrimSpace(req.MessageText)) == 0 {
		return errors.New("message text cannot be empty")
	}

	if len(req.MessageText) > 10000 {
		return errors.New("message text too long (max 10000 characters)")
	}

	// Проверяем тип сообщения
	validTypes := map[string]bool{
		"text":   true,
		"system": true,
		"file":   true,
	}

	if req.MessageType != "" && !validTypes[req.MessageType] {
		return errors.New("invalid message type. Allowed: text, system, file")
	}

	return nil
}

// getMessageType возвращает тип сообщения или дефолтный
func getMessageType(msgType string) string {
	if msgType == "" {
		return "text"
	}
	return msgType
}
