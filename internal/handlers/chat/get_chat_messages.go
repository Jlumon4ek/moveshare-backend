package chat

import (
	"log"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetChatMessages godoc
// @Summary      Get chat messages
// @Description  Returns messages for a specific chat with pagination
// @Tags         Chat
// @Security     BearerAuth
// @Param        chatId path      int     true  "Chat ID"
// @Param        limit  query     int     false "Limit number of messages returned" default(30)
// @Param        offset query     int     false "Offset for pagination" default(0)
// @Param        order  query     string  false "Order of messages" Enums(asc, desc) default(desc)
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chats/{chatId}/messages [get]
func GetChatMessages(chatService service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("GetChatMessages: Starting request")
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			log.Printf("GetChatMessages: Failed to get user ID from context: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		log.Printf("GetChatMessages: User ID: %d", userID)

		chatIDStr := c.Param("chatId")
		log.Printf("GetChatMessages: Chat ID string: %s", chatIDStr)
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil || chatID <= 0 {
			log.Printf("GetChatMessages: Invalid chat ID: %s, error: %v", chatIDStr, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
			return
		}
		log.Printf("GetChatMessages: Chat ID: %d", chatID)

		limitStr := c.DefaultQuery("limit", "30")
		offsetStr := c.DefaultQuery("offset", "0")
		order := c.DefaultQuery("order", "desc")
		log.Printf("GetChatMessages: Parameters - limit: %s, offset: %s, order: %s", limitStr, offsetStr, order)

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 || limit > 100 {
			log.Printf("GetChatMessages: Invalid limit parameter: %s, error: %v", limitStr, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter (0-100)"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}

		if order != "asc" && order != "desc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameter. Use 'asc' or 'desc'"})
			return
		}

		isParticipant, err := chatService.IsUserParticipant(c.Request.Context(), chatID, userID)
		if err != nil {
			log.Printf("GetChatMessages: Failed to verify chat access for user %d, chat %d: %v", userID, chatID, err)
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

		messages, total, err := chatService.GetChatMessages(c.Request.Context(), chatID, userID, limit, offset, order)
		if err != nil {
			log.Printf("GetChatMessages: Failed to get chat messages for user %d, chat %d: %v", userID, chatID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get chat messages",
				"details": err.Error(),
			})
			return
		}

		// Messages are now marked as read via separate API endpoint

		hasNext := (offset + limit) < total
		hasPrev := offset > 0

		c.JSON(http.StatusOK, gin.H{
			"messages": messages,
			"pagination": gin.H{
				"limit":    limit,
				"offset":   offset,
				"total":    total,
				"has_next": hasNext,
				"has_prev": hasPrev,
				"order":    order,
			},
			"chat_info": gin.H{
				"id":             chatID,
				"total_messages": total,
			},
		})
	}
}
