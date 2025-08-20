package notifications

import (
	"moveshare/internal/utils"
	"moveshare/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestNotification godoc
// @Summary      Send test notification
// @Description  Sends a test notification to the current user (for development)
// @Tags         Notifications
// @Security     BearerAuth
// @Param        type query string false "Notification type" Enums(message,job,system) default(system)
// @Param        message query string false "Notification message" default(Test notification)
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /notifications/test [post]
func TestNotification(notificationHub *websocket.NotificationHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		notificationType := c.DefaultQuery("type", "system")
		message := c.DefaultQuery("message", "Test notification")

		var data map[string]interface{}

		switch notificationType {
		case "message":
			data = map[string]interface{}{
				"chat_id":     123,
				"sender_name": "Test User",
				"message":     message,
				"preview":     message,
				"action":      "open_chat",
				"action_url":  "/chats/123",
				"dismissable": true,
				"priority":    "normal",
				"category":    "message",
			}
		case "job":
			data = map[string]interface{}{
				"job_id":      456,
				"status":      "claimed",
				"message":     message,
				"action":      "open_job",
				"action_url":  "/jobs/456",
				"dismissable": true,
				"priority":    "high",
				"category":    "job",
			}
		default: // system
			data = map[string]interface{}{
				"message":     message,
				"level":       "info",
				"action":      "none",
				"action_url":  "",
				"dismissable": true,
				"priority":    "normal",
				"category":    "system",
			}
		}

		notificationHub.SendNotificationToUser(userID, notificationType, data)

		c.JSON(http.StatusOK, gin.H{
			"message": "Test notification sent",
			"type":    notificationType,
			"data":    data,
		})
	}
}
