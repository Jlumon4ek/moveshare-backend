package session

import (
	"log"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetActiveSessions godoc
// @Summary      Get user's active sessions
// @Description  Returns all active sessions for the current user
// @Tags         Session
// @Produce      json
// @Success      200 {object} models.ActiveSessionsResponse "active sessions"
// @Failure      401 {object} models.ErrorResponse "Unauthorized"
// @Failure      500 {object} models.ErrorResponse "Server error"
// @Router       /user/active-sessions [get]
// @Security     BearerAuth
func GetActiveSessions(sessionService service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		sessions, err := sessionService.GetUserActiveSessions(c.Request.Context(), userID)
		if err != nil {
			log.Printf("Error getting active sessions for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get active sessions"})
			return
		}

		c.JSON(http.StatusOK, sessions)
	}
}