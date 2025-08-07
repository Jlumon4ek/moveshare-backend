// internal/models/chat_model.go
package models

import "time"

type ChatConversation struct {
	ID           int64     `json:"id"`
	JobID        int64     `json:"job_id"`
	ClientID     int64     `json:"client_id"`
	ContractorID int64     `json:"contractor_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ChatMessage struct {
	ID             int64      `json:"id"`
	ConversationID int64      `json:"conversation_id"`
	SenderID       int64      `json:"sender_id"`
	MessageText    string     `json:"message_text"`
	MessageType    string     `json:"message_type"`
	IsRead         bool       `json:"is_read"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ChatListItem struct {
	ID              int64     `json:"id"`
	JobID           int64     `json:"job_id"`
	JobTitle        string    `json:"job_title"`
	OtherUserID     int64     `json:"other_user_id"`
	OtherUserName   string    `json:"other_user_name"`
	OtherUserRole   string    `json:"other_user_role"` // client, contractor
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
	LastMessageType string    `json:"last_message_type"`
	UnreadCount     int       `json:"unread_count"`
	IsLastMsgFromMe bool      `json:"is_last_msg_from_me"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ChatMessageResponse struct {
	ID          int64      `json:"id"`
	SenderID    int64      `json:"sender_id"`
	SenderName  string     `json:"sender_name"`
	MessageText string     `json:"message_text"`
	MessageType string     `json:"message_type"`
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	IsFromMe    bool       `json:"is_from_me"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
