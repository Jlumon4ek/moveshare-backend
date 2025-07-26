package admin

import (
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetChatConversationCount handles getting chat conversation count
// @Summary Get total chat conversation count
// @Description Gets the total number of chat conversations in the system
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} map[string]string
// @Router /admin/conversations/count [get]
func GetChatConversationCount(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		count, err := adminService.GetChatConversationCount(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat conversation count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"conversation_count": count})
	}
}
