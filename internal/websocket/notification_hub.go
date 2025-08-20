package websocket

import (
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type NotificationMessage struct {
	Type   string      `json:"type"` // "job_update", "message", "system", "ping", "pong"
	Data   interface{} `json:"data"`
	UserID int64       `json:"user_id"`
	Time   time.Time   `json:"time"`
}

type NotificationClient struct {
	ID     string
	UserID int64
	Conn   *websocket.Conn
	Send   chan NotificationMessage
	Hub    *NotificationHub
}

type NotificationHub struct {
	// Клиенты по пользователям
	Clients map[int64][]*NotificationClient

	// Регистрация клиентов
	Register chan *NotificationClient

	// Отмена регистрации
	Unregister chan *NotificationClient

	// Уведомления для отправки
	Notifications chan NotificationMessage
}

func NewNotificationHub() *NotificationHub {
	return &NotificationHub{
		Clients:       make(map[int64][]*NotificationClient),
		Register:      make(chan *NotificationClient),
		Unregister:    make(chan *NotificationClient),
		Notifications: make(chan NotificationMessage),
	}
}

func (h *NotificationHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = append(h.Clients[client.UserID], client)
			log.Printf("Notification client %s registered for user %d", client.ID, client.UserID)

			// Отправляем подтверждение подключения
			client.Send <- NotificationMessage{
				Type:   "connected",
				Data:   map[string]interface{}{"message": "Successfully connected to notifications"},
				UserID: client.UserID,
				Time:   time.Now(),
			}

		case client := <-h.Unregister:
			if clients, ok := h.Clients[client.UserID]; ok {
				// Находим и удаляем клиента
				for i, c := range clients {
					if c == client {
						h.Clients[client.UserID] = append(clients[:i], clients[i+1:]...)
						close(client.Send)
						break
					}
				}
				// Если у пользователя больше нет подключений, удаляем его из мапы
				if len(h.Clients[client.UserID]) == 0 {
					delete(h.Clients, client.UserID)
				}
			}
			log.Printf("Notification client %s unregistered for user %d", client.ID, client.UserID)

		case notification := <-h.Notifications:
			log.Printf("NotificationHub: Processing notification for user %d (type: %s)", notification.UserID, notification.Type)
			// Отправляем уведомление всем клиентам пользователя
			if clients, ok := h.Clients[notification.UserID]; ok {
				log.Printf("NotificationHub: Found %d clients for user %d", len(clients), notification.UserID)
				for i, client := range clients {
					log.Printf("NotificationHub: Sending to client %d for user %d", i, notification.UserID)
					select {
					case client.Send <- notification:
						log.Printf("NotificationHub: Successfully sent notification to user %d client %d", notification.UserID, i)
					default:
						// Клиент не отвечает, удаляем его
						log.Printf("NotificationHub: Client %d not responding, removing for user %d", i, notification.UserID)
						h.removeClient(client)
					}
				}
			} else {
				log.Printf("NotificationHub: No clients found for user %d", notification.UserID)
			}
		}
	}
}

func (h *NotificationHub) removeClient(client *NotificationClient) {
	if clients, ok := h.Clients[client.UserID]; ok {
		for i, c := range clients {
			if c == client {
				h.Clients[client.UserID] = append(clients[:i], clients[i+1:]...)
				close(client.Send)
				break
			}
		}
		if len(h.Clients[client.UserID]) == 0 {
			delete(h.Clients, client.UserID)
		}
	}
}

func (c *NotificationClient) ReadPump() {
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
		var msg NotificationMessage
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

func (c *NotificationClient) WritePump() {
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

func (c *NotificationClient) handleIncomingMessage(msg NotificationMessage) {
	switch msg.Type {
	case "ping":
		// Отвечаем на ping
		c.Send <- NotificationMessage{
			Type:   "pong",
			UserID: c.UserID,
			Time:   time.Now(),
		}
	}
}

// SendNotificationToUser отправляет уведомление конкретному пользователю
func (h *NotificationHub) SendNotificationToUser(userID int64, notificationType string, data interface{}) {
	notification := NotificationMessage{
		Type:   notificationType,
		Data:   data,
		UserID: userID,
		Time:   time.Now(),
	}

	log.Printf("NotificationHub: Queueing notification type '%s' to user %d", notificationType, userID)
	log.Printf("NotificationHub: Current connected users: %v", h.getConnectedUsers())
	h.Notifications <- notification
}

// getConnectedUsers возвращает список подключенных пользователей для отладки
func (h *NotificationHub) getConnectedUsers() []int64 {
	var users []int64
	for userID, clients := range h.Clients {
		if len(clients) > 0 {
			users = append(users, userID)
		}
	}
	return users
}

func GenerateNotificationClientID(userID int64) string {
	return "notification_" + strconv.FormatInt(userID, 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}