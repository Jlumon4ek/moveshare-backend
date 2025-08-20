package notifications

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetNotifications godoc
// @Summary      Get user notifications
// @Description  Returns a paginated list of notifications for the authenticated user
// @Tags         Notifications
// @Security     BearerAuth
// @Param        limit    query     int     false "Limit number of notifications returned" default(20)
// @Param        offset   query     int     false "Offset for pagination" default(0)
// @Param        type     query     string  false "Filter by notification type" Enums(job_application, job_update, payment, document_upload, new_job, review, message, system)
// @Param        unread   query     bool    false "Show only unread notifications" default(false)
// @Produce      json
// @Success      200  {object}  models.NotificationListResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notifications [get]
func GetNotifications(notificationService service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Parse query parameters
		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")
		typeFilter := c.DefaultQuery("type", "all")
		unreadOnlyStr := c.DefaultQuery("unread", "false")

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

		unreadOnly, err := strconv.ParseBool(unreadOnlyStr)
		if err != nil {
			unreadOnly = false
		}

		// Get notifications
		response, err := notificationService.GetUserNotifications(c.Request.Context(), userID, limit, offset, typeFilter, unreadOnly)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get notifications",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetNotificationStats godoc
// @Summary      Get notification statistics
// @Description  Returns notification statistics for the authenticated user
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.NotificationStats
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notifications/stats [get]
func GetNotificationStats(notificationService service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		stats, err := notificationService.GetNotificationStats(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get notification stats",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}