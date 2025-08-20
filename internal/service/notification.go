package service

import (
	"context"
	"fmt"
	"log"
	"moveshare/internal/models"
	"moveshare/internal/repository/notifications"
	"moveshare/internal/websocket"
	"time"
)

type NotificationService interface {
	// Core CRUD operations
	CreateNotification(ctx context.Context, req *models.NotificationRequest) (*models.Notification, error)
	GetUserNotifications(ctx context.Context, userID int64, limit, offset int, typeFilter string, unreadOnly bool) (*models.NotificationListResponse, error)
	GetNotificationByID(ctx context.Context, id, userID int64) (*models.Notification, error)
	MarkAsRead(ctx context.Context, id, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	DeleteNotification(ctx context.Context, id, userID int64) error
	DeleteAllNotifications(ctx context.Context, userID int64) error
	GetNotificationStats(ctx context.Context, userID int64) (*models.NotificationStats, error)

	// Job-related notifications  
	NotifyJobApplication(ctx context.Context, jobOwnerID, applicantID, jobID int64, applicantName string) error
	NotifyJobClaimed(ctx context.Context, jobOwnerID, contractorID, jobID int64, contractorName, jobTitle string) error
	NotifyJobCompleted(ctx context.Context, jobOwnerID, contractorID, jobID int64, jobTitle string) error
	NotifyDocumentUploaded(ctx context.Context, recipientID, uploaderID, jobID int64, uploaderName, documentType string) error
	NotifyPaymentRequired(ctx context.Context, userID, jobID int64, amount float64, dueDate time.Time) error
	NotifyNewReview(ctx context.Context, userID, reviewerID, jobID int64, reviewerName string, rating int) error
	
	// System notifications
	NotifyNewMatchingJob(ctx context.Context, userID, jobID int64, jobTitle, route string, estimatedPay float64) error
	NotifySystemAnnouncement(ctx context.Context, userID int64, title, message string, priority models.NotificationPriority) error
	
	// WebSocket real-time notifications
	NotifyJobUpdate(userID int64, jobID int64, status string, message string)
	NotifyNewMessage(userID int64, chatID int64, senderName string, messageText string)
	NotifyUnreadCountChange(userID int64, newUnreadCount int)
	NotifySystemMessage(userID int64, message string, level string)
	
	// Cleanup
	CleanupExpiredNotifications(ctx context.Context) (int, error)
}

type notificationService struct {
	repo notifications.NotificationRepository
	hub  *websocket.NotificationHub
}

func NewNotificationService(hub *websocket.NotificationHub, repo notifications.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
		hub:  hub,
	}
}

// Core CRUD operations

func (s *notificationService) CreateNotification(ctx context.Context, req *models.NotificationRequest) (*models.Notification, error) {
	return s.repo.Create(ctx, req)
}

func (s *notificationService) GetUserNotifications(ctx context.Context, userID int64, limit, offset int, typeFilter string, unreadOnly bool) (*models.NotificationListResponse, error) {
	notifications, total, err := s.repo.GetByUserID(ctx, userID, limit, offset, typeFilter, unreadOnly)
	if err != nil {
		return nil, err
	}

	response := &models.NotificationListResponse{
		Notifications: notifications,
		Total:        total,
	}
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset
	response.Pagination.HasNext = offset+len(notifications) < total
	response.Pagination.HasPrev = offset > 0
	
	return response, nil
}

func (s *notificationService) GetNotificationByID(ctx context.Context, id, userID int64) (*models.Notification, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id, userID int64) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) DeleteNotification(ctx context.Context, id, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *notificationService) DeleteAllNotifications(ctx context.Context, userID int64) error {
	return s.repo.DeleteAll(ctx, userID)
}

func (s *notificationService) GetNotificationStats(ctx context.Context, userID int64) (*models.NotificationStats, error) {
	return s.repo.GetStats(ctx, userID)
}

func (s *notificationService) CleanupExpiredNotifications(ctx context.Context) (int, error) {
	return s.repo.CleanupExpired(ctx)
}

// WebSocket real-time notifications

func (s *notificationService) NotifyJobUpdate(userID int64, jobID int64, status string, message string) {
	var actionURL string
	var priority string

	switch status {
	case "claimed":
		actionURL = fmt.Sprintf("/jobs/%d/details", jobID)
		priority = "high"
	case "completed":
		actionURL = fmt.Sprintf("/jobs/%d/review", jobID)
		priority = "high"
	case "pending":
		actionURL = fmt.Sprintf("/jobs/%d", jobID)
		priority = "normal"
	default:
		actionURL = fmt.Sprintf("/jobs/%d", jobID)
		priority = "normal"
	}

	data := map[string]interface{}{
		"job_id":     jobID,
		"status":     status,
		"message":    message,
		"action":     "view_job",
		"action_url": actionURL,
		"priority":   priority,
		"category":   "job_update",
	}

	log.Printf("Sending job update notification to user %d: job %d status %s", userID, jobID, status)
	s.hub.SendNotificationToUser(userID, "job_update", data)
}

func (s *notificationService) NotifyNewMessage(userID int64, chatID int64, senderName string, messageText string) {
	preview := messageText
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}

	data := map[string]interface{}{
		"chat_id":     chatID,
		"sender_name": senderName,
		"preview":     preview,
		"action":      "open_chat",
		"action_url":  fmt.Sprintf("/chats?chat=%d", chatID),
		"priority":    "normal",
		"category":    "new_message",
	}

	log.Printf("NotificationService: Notification data: %+v", data)
	s.hub.SendNotificationToUser(userID, "new_message", data)
}

func (s *notificationService) NotifyUnreadCountChange(userID int64, newUnreadCount int) {
	data := map[string]interface{}{
		"unread_count": newUnreadCount,
		"category":     "chat_count",
	}

	log.Printf("NotificationService: Sending unread count update to user %d: %d unread", userID, newUnreadCount)
	s.hub.SendNotificationToUser(userID, "unread_count_update", data)
}

func (s *notificationService) NotifySystemMessage(userID int64, message string, level string) {
	var priority string
	var dismissable bool

	switch level {
	case "error":
		priority = "high"
		dismissable = false
	case "warning":
		priority = "normal"
		dismissable = true
	case "success":
		priority = "normal"
		dismissable = true
	case "info":
		priority = "low"
		dismissable = true
	default:
		priority = "normal"
		dismissable = true
	}

	data := map[string]interface{}{
		"message":     message,
		"level":       level,
		"action":      "none",
		"action_url":  "",
		"dismissable": dismissable,
		"priority":    priority,
		"category":    "system",
	}

	log.Printf("Sending system notification to user %d: %s (%s)", userID, message, level)
	s.hub.SendNotificationToUser(userID, "system", data)
}

// Job-related notification implementations

func (s *notificationService) NotifyJobApplication(ctx context.Context, jobOwnerID, applicantID, jobID int64, applicantName string) error {
	req := &models.NotificationRequest{
		UserID:        jobOwnerID,
		Type:          models.NotificationTypeJobApplication,
		Title:         "New Job Application",
		Message:       fmt.Sprintf("%s has applied for your job. Review their profile and accept or decline the application.", applicantName),
		JobID:         &jobID,
		RelatedUserID: &applicantID,
		Priority:      models.NotificationPriorityHigh,
		Actions: []models.NotificationAction{
			{Label: "View Application", Action: "view_job", URL: fmt.Sprintf("/jobs/%d", jobID), Primary: true},
			{Label: "View Profile", Action: "view_profile", URL: fmt.Sprintf("/profile/%d", applicantID)},
			{Label: "Mark as Read", Action: "mark_read"},
		},
		Metadata: map[string]interface{}{
			"applicant_name": applicantName,
			"applicant_id":   applicantID,
		},
	}

	_, err := s.repo.Create(ctx, req)
	return err
}

func (s *notificationService) NotifyJobClaimed(ctx context.Context, jobOwnerID, contractorID, jobID int64, contractorName, jobTitle string) error {
	req := &models.NotificationRequest{
		UserID:        jobOwnerID,
		Type:          models.NotificationTypeJobClaimed,
		Title:         "Job Claimed Successfully",
		Message:       fmt.Sprintf("Your job '%s' has been claimed by %s. You can now coordinate the details via chat.", jobTitle, contractorName),
		JobID:         &jobID,
		RelatedUserID: &contractorID,
		Priority:      models.NotificationPriorityHigh,
		Actions: []models.NotificationAction{
			{Label: "View Job", Action: "view_job", URL: fmt.Sprintf("/jobs/%d", jobID), Primary: true},
			{Label: "Open Chat", Action: "open_chat", URL: fmt.Sprintf("/chats?job=%d", jobID)},
			{Label: "Mark as Read", Action: "mark_read"},
		},
		Metadata: map[string]interface{}{
			"contractor_name": contractorName,
			"contractor_id":   contractorID,
			"job_title":       jobTitle,
		},
	}

	_, err := s.repo.Create(ctx, req)
	return err
}

func (s *notificationService) NotifyJobCompleted(ctx context.Context, jobOwnerID, contractorID, jobID int64, jobTitle string) error {
	req := &models.NotificationRequest{
		UserID:        jobOwnerID,
		Type:          models.NotificationTypeJobCompleted,
		Title:         "Job Completed",
		Message:       fmt.Sprintf("The job '%s' has been marked as completed. Please review the work and leave feedback.", jobTitle),
		JobID:         &jobID,
		RelatedUserID: &contractorID,
		Priority:      models.NotificationPriorityHigh,
		Actions: []models.NotificationAction{
			{Label: "Review Work", Action: "review_job", URL: fmt.Sprintf("/jobs/%d/review", jobID), Primary: true},
			{Label: "View Job", Action: "view_job", URL: fmt.Sprintf("/jobs/%d", jobID)},
			{Label: "Mark as Read", Action: "mark_read"},
		},
		Metadata: map[string]interface{}{
			"contractor_id": contractorID,
			"job_title":     jobTitle,
		},
	}

	_, err := s.repo.Create(ctx, req)
	return err
}

func (s *notificationService) NotifyDocumentUploaded(ctx context.Context, recipientID, uploaderID, jobID int64, uploaderName, documentType string) error {
	req := &models.NotificationRequest{
		UserID:        recipientID,
		Type:          models.NotificationTypeDocumentUpload,
		Title:         "New Document Uploaded",
		Message:       fmt.Sprintf("%s has uploaded a %s document for the job. Please review the document.", uploaderName, documentType),
		JobID:         &jobID,
		RelatedUserID: &uploaderID,
		Priority:      models.NotificationPriorityNormal,
		Actions: []models.NotificationAction{
			{Label: "View Document", Action: "view_document", URL: fmt.Sprintf("/jobs/%d/documents", jobID), Primary: true},
			{Label: "View Job", Action: "view_job", URL: fmt.Sprintf("/jobs/%d", jobID)},
			{Label: "Mark as Read", Action: "mark_read"},
		},
		Metadata: map[string]interface{}{
			"uploader_name":  uploaderName,
			"uploader_id":    uploaderID,
			"document_type":  documentType,
		},
	}

	_, err := s.repo.Create(ctx, req)
	return err
}

func (s *notificationService) NotifyPaymentRequired(ctx context.Context, userID, jobID int64, amount float64, dueDate time.Time) error {
	req := &models.NotificationRequest{
		UserID:   userID,
		Type:     models.NotificationTypePayment,
		Title:    "Payment Required",
		Message:  fmt.Sprintf("Payment of $%.2f is required for your job. Please complete payment by %s to avoid service interruption.", amount, dueDate.Format("Jan 2, 2006")),
		JobID:    &jobID,
		Priority: models.NotificationPriorityHigh,
		Actions: []models.NotificationAction{
			{Label: "Pay Now", Action: "pay_now", URL: fmt.Sprintf("/jobs/%d/payment", jobID), Primary: true},
			{Label: "View Job", Action: "view_job", URL: fmt.Sprintf("/jobs/%d", jobID)},
			{Label: "Remind Me Later", Action: "remind_later"},
		},
		Metadata: map[string]interface{}{
			"amount":   amount,
			"due_date": dueDate.Format(time.RFC3339),
		},
		ExpiresAt: &dueDate,
	}

	_, err := s.repo.Create(ctx, req)
	return err
}

func (s *notificationService) NotifyNewReview(ctx context.Context, userID, reviewerID, jobID int64, reviewerName string, rating int) error {
	req := &models.NotificationRequest{
		UserID:        userID,
		Type:          models.NotificationTypeReview,
		Title:         "New Review Received",
		Message:       fmt.Sprintf("%s has left you a %d-star review for your recent job. Check out what they said!", reviewerName, rating),
		JobID:         &jobID,
		RelatedUserID: &reviewerID,
		Priority:      models.NotificationPriorityNormal,
		Actions: []models.NotificationAction{
			{Label: "View Review", Action: "view_review", URL: fmt.Sprintf("/jobs/%d/review", jobID), Primary: true},
			{Label: "View Job", Action: "view_job", URL: fmt.Sprintf("/jobs/%d", jobID)},
			{Label: "Mark as Read", Action: "mark_read"},
		},
		Metadata: map[string]interface{}{
			"reviewer_name": reviewerName,
			"reviewer_id":   reviewerID,
			"rating":        rating,
		},
	}

	_, err := s.repo.Create(ctx, req)
	return err
}

func (s *notificationService) NotifyNewMatchingJob(ctx context.Context, userID, jobID int64, jobTitle, route string, estimatedPay float64) error {
	req := &models.NotificationRequest{
		UserID:   userID,
		Type:     models.NotificationTypeNewJob,
		Title:    "New Job Matching Your Route",
		Message:  fmt.Sprintf("A new job '%s' on route %s has been posted. Estimated payout: $%.2f", jobTitle, route, estimatedPay),
		JobID:    &jobID,
		Priority: models.NotificationPriorityNormal,
		Actions: []models.NotificationAction{
			{Label: "View Job", Action: "view_job", URL: fmt.Sprintf("/jobs/%d", jobID), Primary: true},
			{Label: "Apply Now", Action: "apply_job", URL: fmt.Sprintf("/jobs/%d/apply", jobID)},
			{Label: "Dismiss", Action: "dismiss"},
		},
		Metadata: map[string]interface{}{
			"job_title":      jobTitle,
			"route":          route,
			"estimated_pay":  estimatedPay,
		},
		ExpiresAt: func() *time.Time { t := time.Now().Add(7 * 24 * time.Hour); return &t }(),
	}

	_, err := s.repo.Create(ctx, req)
	return err
}

func (s *notificationService) NotifySystemAnnouncement(ctx context.Context, userID int64, title, message string, priority models.NotificationPriority) error {
	req := &models.NotificationRequest{
		UserID:   userID,
		Type:     models.NotificationTypeSystem,
		Title:    title,
		Message:  message,
		Priority: priority,
		Actions: []models.NotificationAction{
			{Label: "Mark as Read", Action: "mark_read"},
		},
		Metadata: map[string]interface{}{
			"source": "system",
		},
	}

	_, err := s.repo.Create(ctx, req)
	return err
}