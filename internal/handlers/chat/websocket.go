package chat

import (
	"context"
	"log"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// В production нужно проверять домены
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSMessage struct {
	Type   string      `json:"type"` // "message", "typing", "read", "error", "ping", "pong"
	Data   interface{} `json:"data"`
	ChatID int64       `json:"chat_id"`
	UserID int64       `json:"user_id"`
	Time   time.Time   `json:"time"`
}

type Client struct {
	ID     string
	UserID int64
	ChatID int64
	Conn   *websocket.Conn
	Send   chan WSMessage
	Hub    *Hub
}

type Hub struct {
	// Клиенты по чатам
	Clients map[int64]map[*Client]bool

	// Регистрация клиентов
	Register chan *Client

	// Отмена регистрации
	Unregister chan *Client

	// Сообщения для броадкаста
	Broadcast chan WSMessage

	// Сервис чата
	ChatService service.ChatService
}

func NewHub(chatService service.ChatService) *Hub {
	return &Hub{
		Clients:     make(map[int64]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan WSMessage),
		ChatService: chatService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// Регистрируем клиента
			if h.Clients[client.ChatID] == nil {
				h.Clients[client.ChatID] = make(map[*Client]bool)
			}
			h.Clients[client.ChatID][client] = true

			log.Printf("Client %s joined chat %d (user %d)", client.ID, client.ChatID, client.UserID)

			// Отправляем подтверждение подключения
			client.Send <- WSMessage{
				Type:   "connected",
				Data:   gin.H{"message": "Successfully connected to chat"},
				ChatID: client.ChatID,
				Time:   time.Now(),
			}

		case client := <-h.Unregister:
			// Удаляем клиента
			if clients, ok := h.Clients[client.ChatID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)

					// Если в чате больше нет клиентов, удаляем чат из хаба
					if len(clients) == 0 {
						delete(h.Clients, client.ChatID)
					}
				}
			}
			log.Printf("Client %s left chat %d", client.ID, client.ChatID)

		case message := <-h.Broadcast:
			// ✅ ИСПРАВЛЕНИЕ: Отправляем сообщение всем клиентам в чате, КРОМЕ отправителя
			if clients, ok := h.Clients[message.ChatID]; ok {
				for client := range clients {
					// ✅ Не отправляем сообщение отправителю
					if client.UserID != message.UserID {
						select {
						case client.Send <- message:
							log.Printf("Sent message to user %d in chat %d", client.UserID, message.ChatID)
						default:
							// Клиент не отвечает, удаляем его
							log.Printf("Client %d not responding, removing from chat %d", client.UserID, message.ChatID)
							delete(clients, client)
							close(client.Send)
						}
					}
				}
			}
		}
	}
}

// WebSocketChat godoc
// @Summary      WebSocket connection for chat
// @Description  Establishes WebSocket connection for real-time chat messaging
// @Tags         Chat
// @Param        chatId path      int     true  "Chat ID"
// @Param        token  query     string  true  "JWT token for authentication"
// @Router       /chats/{chatId}/ws [get]
func WebSocketChat(hub *Hub, jwtAuth service.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем chatId
		chatIDStr := c.Param("chatId")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil || chatID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
			return
		}

		// Аутентификация через query parameter (для WebSocket)
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

		// Проверяем доступ к чату
		isParticipant, err := hub.ChatService.IsUserParticipant(c.Request.Context(), chatID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify chat access"})
			return
		}

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// Апгрейдим соединение до WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		// Создаем клиента
		client := &Client{
			ID:     generateClientID(userID, chatID),
			UserID: userID,
			ChatID: chatID,
			Conn:   conn,
			Send:   make(chan WSMessage, 256),
			Hub:    hub,
		}

		// Регистрируем клиента в хабе
		hub.Register <- client

		// Запускаем горутины для чтения и записи
		go client.writePump()
		go client.readPump()
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Обрабатываем входящие сообщения
		c.handleIncomingMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("WriteJSON error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleIncomingMessage(msg WSMessage) {
	switch msg.Type {
	case "ping":
		// Отвечаем на ping
		c.Send <- WSMessage{
			Type:   "pong",
			ChatID: c.ChatID,
			Time:   time.Now(),
		}

	case "typing":
		// Отправляем уведомление о печати другим участникам
		c.Hub.Broadcast <- WSMessage{
			Type:   "typing",
			Data:   gin.H{"user_id": c.UserID, "is_typing": true},
			ChatID: c.ChatID,
			UserID: c.UserID,
			Time:   time.Now(),
		}

	case "stop_typing":
		// Отправляем уведомление о прекращении печати
		c.Hub.Broadcast <- WSMessage{
			Type:   "typing",
			Data:   gin.H{"user_id": c.UserID, "is_typing": false},
			ChatID: c.ChatID,
			UserID: c.UserID,
			Time:   time.Now(),
		}

	case "read":
		// Помечаем сообщения как прочитанные
		go func() {
			ctx := context.Background()
			err := c.Hub.ChatService.MarkMessagesAsRead(ctx, c.ChatID, c.UserID)
			if err != nil {
				log.Printf("Failed to mark messages as read: %v", err)
			}
		}()
	}
}

// ✅ ИСПРАВЛЕНИЕ: BroadcastMessage отправляет сообщение всем участникам чата КРОМЕ отправителя
func (h *Hub) BroadcastMessage(chatID, senderID int64, message *models.ChatMessageResponse) {
	wsMessage := WSMessage{
		Type:   "message",
		Data:   message,
		ChatID: chatID,
		UserID: senderID, // Это ID отправителя, чтобы не отправлять ему же
		Time:   time.Now(),
	}

	log.Printf("Broadcasting message from user %d to chat %d", senderID, chatID)
	h.Broadcast <- wsMessage
}

func generateClientID(userID, chatID int64) string {
	return strconv.FormatInt(userID, 10) + "_" + strconv.FormatInt(chatID, 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
