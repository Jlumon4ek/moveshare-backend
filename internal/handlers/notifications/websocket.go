package notifications

import (
	"log"
	"moveshare/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// В production нужно проверять домены
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// JWTAuthenticator интерфейс для проверки JWT токенов
type JWTAuthenticator interface {
	ValidateToken(token string) (int64, error)
}

// WebSocketNotifications godoc
// @Summary      WebSocket connection for notifications
// @Description  Establishes WebSocket connection for real-time notifications
// @Tags         Notifications
// @Param        token  query     string  true  "JWT token for authentication"
// @Router       /notifications/ws [get]
func WebSocketNotifications(hub *websocket.NotificationHub, jwtAuth JWTAuthenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Аутентификация через query parameter
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}

		userID, err := jwtAuth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Апгрейдим соединение до WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		// Создаем клиента
		client := &websocket.NotificationClient{
			ID:     websocket.GenerateNotificationClientID(userID),
			UserID: userID,
			Conn:   conn,
			Send:   make(chan websocket.NotificationMessage, 256),
			Hub:    hub,
		}

		// Регистрируем клиента в хабе
		hub.Register <- client

		// Запускаем горутины для чтения и записи
		go client.WritePump()
		go client.ReadPump()
	}
}

