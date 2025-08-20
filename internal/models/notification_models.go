package models

import (
	"encoding/json"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeJobApplication NotificationType = "job_application" // Someone applied to your job
	NotificationTypeJobUpdate      NotificationType = "job_update"      // Job status changed
	NotificationTypeJobClaimed     NotificationType = "job_claimed"     // Your job was claimed
	NotificationTypeJobCompleted   NotificationType = "job_completed"   // Job was completed
	NotificationTypePayment        NotificationType = "payment"         // Payment related
	NotificationTypeDocumentUpload NotificationType = "document_upload" // Document uploaded
	NotificationTypeNewJob         NotificationType = "new_job"         // New job matching criteria
	NotificationTypeReview         NotificationType = "review"          // New review received
	NotificationTypeMessage        NotificationType = "message"         // New chat message
	NotificationTypeSystem         NotificationType = "system"          // System announcements
)

// NotificationPriority represents the priority of notification
type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityNormal NotificationPriority = "normal"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// NotificationAction represents an action that can be taken on a notification
type NotificationAction struct {
	Label   string `json:"label"`
	Action  string `json:"action"`  // action type: 'view_job', 'open_chat', 'mark_read', etc.
	URL     string `json:"url"`     // URL to navigate to
	Primary bool   `json:"primary"` // is this the primary action
}

// Notification represents a notification in the database
type Notification struct {
	ID             int64                  `json:"id" db:"id"`
	UserID         int64                  `json:"user_id" db:"user_id"`
	Type           NotificationType       `json:"type" db:"type"`
	Title          string                 `json:"title" db:"title"`
	Message        string                 `json:"message" db:"message"`
	JobID          *int64                 `json:"job_id,omitempty" db:"job_id"`
	ChatID         *int64                 `json:"chat_id,omitempty" db:"chat_id"`
	RelatedUserID  *int64                 `json:"related_user_id,omitempty" db:"related_user_id"`
	IsRead         bool                   `json:"is_read" db:"is_read"`
	Priority       NotificationPriority   `json:"priority" db:"priority"`
	Actions        []NotificationAction   `json:"actions" db:"actions"`
	Metadata       map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	ReadAt         *time.Time             `json:"read_at,omitempty" db:"read_at"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty" db:"expires_at"`
}

// NotificationRequest represents a request to create a notification
type NotificationRequest struct {
	UserID        int64                  `json:"user_id"`
	Type          NotificationType       `json:"type"`
	Title         string                 `json:"title"`
	Message       string                 `json:"message"`
	JobID         *int64                 `json:"job_id,omitempty"`
	ChatID        *int64                 `json:"chat_id,omitempty"`
	RelatedUserID *int64                 `json:"related_user_id,omitempty"`
	Priority      NotificationPriority   `json:"priority"`
	Actions       []NotificationAction   `json:"actions"`
	Metadata      map[string]interface{} `json:"metadata"`
	ExpiresAt     *time.Time             `json:"expires_at,omitempty"`
}

// NotificationListResponse represents a paginated list of notifications
type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	Total         int            `json:"total"`
	Unread        int            `json:"unread"`
	Pagination    struct {
		Limit    int  `json:"limit"`
		Offset   int  `json:"offset"`
		HasNext  bool `json:"has_next"`
		HasPrev  bool `json:"has_prev"`
	} `json:"pagination"`
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	Total  int `json:"total"`
	Unread int `json:"unread"`
	ByType map[NotificationType]int `json:"by_type"`
}

// Scan implements the sql.Scanner interface for Actions field
func (n *Notification) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		return json.Unmarshal(s, &n.Actions)
	case string:
		return json.Unmarshal([]byte(s), &n.Actions)
	}
	return nil
}