package chat

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
	GetUserChats(ctx context.Context, userID int64, limit, offset int) ([]models.ChatListItem, int, error)
	GetChatMessages(ctx context.Context, chatID, userID int64, limit, offset int, order string) ([]models.ChatMessageResponse, int, error)
	IsUserParticipant(ctx context.Context, chatID, userID int64) (bool, error)
	MarkMessagesAsRead(ctx context.Context, chatID, userID int64) error
	SendMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessageResponse, error)
	IsChatActive(ctx context.Context, chatID int64) (bool, error)
	UpdateChatActivity(ctx context.Context, chatID int64) error
	CreateChat(ctx context.Context, jobID, clientID, contractorID int64) (int64, error)
	FindExistingChat(ctx context.Context, jobID, client_id, contractor_id int64) (int64, error)
	HasJobAccess(ctx context.Context, jobID, userID1, userID2 int64) (bool, error)
	GetChatParticipants(ctx context.Context, chatID int64) ([]models.ChatParticipant, error)
	GetUserUnreadCount(ctx context.Context, userID int64) (int, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) ChatRepository {
	return &repository{db: db}
}
