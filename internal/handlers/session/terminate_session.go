package session

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TerminateSession godoc
// @Summary      Terminate a specific session
// @Description  Terminates a specific session by ID
// @Tags         Session
// @Param        session_id path int true "Session ID"
// @Produce      json
// @Success      200 {object} map[string]string "success message"
// @Failure      400 {object} models.ErrorResponse "Bad request"
// @Failure      401 {object} models.ErrorResponse "Unauthorized"
// @Failure      500 {object} models.ErrorResponse "Server error"
// @Router       /user/sessions/{session_id}/terminate [delete]
// @Security     BearerAuth
func TerminateSession(sessionService service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		sessionIDStr := c.Param("session_id")
		sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid session ID"})
			return
		}

		err = sessionService.TerminateSession(c.Request.Context(), sessionID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to terminate session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Session terminated successfully",
		})
	}
}

// TerminateAllSessions godoc
// @Summary      Terminate all sessions except current
// @Description  Terminates all user sessions except the current one
// @Tags         Session
// @Produce      json
// @Success      200 {object} map[string]string "success message"
// @Failure      401 {object} models.ErrorResponse "Unauthorized"
// @Failure      500 {object} models.ErrorResponse "Server error"
// @Router       /user/sessions/terminate-all [delete]
// @Security     BearerAuth
func TerminateAllSessions(sessionService service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		err = sessionService.TerminateAllSessions(c.Request.Context(), userID, true) // Keep current session
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to terminate sessions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "All other sessions terminated successfully",
		})
	}
}