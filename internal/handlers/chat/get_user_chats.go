// internal/handlers/chat/get_user_chats.go
package chat

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserChats godoc
// @Summary      Get user chats
// @Description  Returns a list of all chats for the authenticated user with preview of last message
// @Tags         Chat
// @Security     BearerAuth
// @Param        limit  query     int     false  "Limit number of chats returned" default(20)
// @Param        offset query     int     false  "Offset for pagination" default(0)
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chats [get]
func GetUserChats(chatService service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 || limit > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter (0-100)"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}

		chats, total, err := chatService.GetUserChats(c.Request.Context(), userID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get user chats",
				"details": err.Error(),
			})
			return
		}

		if len(chats) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"chats":   []interface{}{},
				"total":   0,
				"message": "No chats found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"chats": chats,
			"total": total,
			"pagination": gin.H{
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}
