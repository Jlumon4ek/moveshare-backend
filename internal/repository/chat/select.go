package chat

import (
	"context"
	"time"
)

type MessagePreview struct {
	OtherUserID int64
	LastMessage string
	SentAt      time.Time
}

func (r *repository) GetChatList(ctx context.Context, userID int64) ([]*MessagePreview, error) {
	return nil, nil
}
