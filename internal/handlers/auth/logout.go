package auth

import (
	"log"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout godoc
// @Summary      Logout user
// @Description  Terminates the current user session
// @Tags         Auth
// @Produce      json
// @Success      200 {object} map[string]string "success message"
// @Failure      401 {object} models.ErrorResponse "Unauthorized"
// @Failure      500 {object} models.ErrorResponse "Server error"
// @Router       /auth/logout [post]
// @Security     BearerAuth
func Logout(sessionService service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		// Get session ID from JWT token (set by middleware)
		sessionID, exists := c.Get("sessionID")
		if exists {
			err = sessionService.TerminateSession(c.Request.Context(), sessionID.(int64), userID)
			if err != nil {
				log.Printf("Error terminating session %d during logout: %v", sessionID, err)
			}
		} else {
			log.Printf("Session ID not found in context during logout for user %d", userID)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Logged out successfully",
		})
	}
}