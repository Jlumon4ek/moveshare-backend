package chat

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
}

type repository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) ChatRepository {
	return &repository{db: db}
}
