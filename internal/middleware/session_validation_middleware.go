package middleware

import (
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SessionValidationMiddleware проверяет, что сессия из JWT токена всё ещё активна
// Используется после AuthMiddleware для критичных операций
func SessionValidationMiddleware(sessionService service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, exists := c.Get("sessionID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID not found"})
			c.Abort()
			return
		}

		// Если sessionID равен 0, это токен без сессии (например, refresh token)
		// Пропускаем проверку сессии для таких токенов
		if sessionID.(int64) == 0 {
			c.Next()
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			c.Abort()
			return
		}

		// Проверяем, что сессия всё ещё существует в базе данных
		sessions, err := sessionService.GetUserActiveSessions(c.Request.Context(), userID.(int64))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session validation failed"})
			c.Abort()
			return
		}

		// Проверяем, что конкретная сессия из JWT всё ещё существует
		sessionExists := false
		for _, session := range sessions.Sessions {
			if session.ID == sessionID.(int64) {
				sessionExists = true
				break
			}
		}

		if !sessionExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session has been terminated"})
			c.Abort()
			return
		}

		c.Next()
	}
}