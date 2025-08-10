package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository/chat"
)

type ChatService interface {
	GetUserChats(ctx context.Context, userID int64, limit, offset int) ([]models.ChatListItem, int, error)
	GetChatMessages(ctx context.Context, chatID, userID int64, limit, offset int, order string) ([]models.ChatMessageResponse, int, error)
	IsUserParticipant(ctx context.Context, chatID, userID int64) (bool, error)
	MarkMessagesAsRead(ctx context.Context, chatID, userID int64) error
	SendMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessageResponse, error)
	IsChatActive(ctx context.Context, chatID int64) (bool, error)
	UpdateChatActivity(ctx context.Context, chatID int64) error
	CreateChat(ctx context.Context, jobID, clientID, contractorID int64) (int64, error)
	FindExistingChat(ctx context.Context, jobID, userID1, userID2 int64) (int64, error)
	HasJobAccess(ctx context.Context, jobID, userID1, userID2 int64) (bool, error)
}

type chatService struct {
	chatRepo chat.ChatRepository
}

func NewChatService(chatRepo chat.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) GetUserChats(ctx context.Context, userID int64, limit, offset int) ([]models.ChatListItem, int, error) {
	return s.chatRepo.GetUserChats(ctx, userID, limit, offset)
}

func (s *chatService) GetChatMessages(ctx context.Context, chatID, userID int64, limit, offset int, order string) ([]models.ChatMessageResponse, int, error) {
	return s.chatRepo.GetChatMessages(ctx, chatID, userID, limit, offset, order)
}

func (s *chatService) IsUserParticipant(ctx context.Context, chatID, userID int64) (bool, error) {
	return s.chatRepo.IsUserParticipant(ctx, chatID, userID)
}

func (s *chatService) MarkMessagesAsRead(ctx context.Context, chatID, userID int64) error {
	return s.chatRepo.MarkMessagesAsRead(ctx, chatID, userID)
}

func (s *chatService) SendMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessageResponse, error) {
	return s.chatRepo.SendMessage(ctx, message)
}

func (s *chatService) IsChatActive(ctx context.Context, chatID int64) (bool, error) {
	return s.chatRepo.IsChatActive(ctx, chatID)
}

func (s *chatService) UpdateChatActivity(ctx context.Context, chatID int64) error {
	return s.chatRepo.UpdateChatActivity(ctx, chatID)
}

func (s *chatService) CreateChat(ctx context.Context, jobID, clientID, contractorID int64) (int64, error) {
	return s.chatRepo.CreateChat(ctx, jobID, clientID, contractorID)
}

func (s *chatService) FindExistingChat(ctx context.Context, jobID, userID1, userID2 int64) (int64, error) {
	return s.chatRepo.FindExistingChat(ctx, jobID, userID1, userID2)
}

func (s *chatService) HasJobAccess(ctx context.Context, jobID, userID1, userID2 int64) (bool, error) {
	return s.chatRepo.HasJobAccess(ctx, jobID, userID1, userID2)
}
