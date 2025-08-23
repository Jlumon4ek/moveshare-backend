package session

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.UserSession) error
	GetUserSessions(ctx context.Context, userID int64) ([]models.UserSession, error)
	GetSessionByToken(ctx context.Context, sessionToken string) (*models.UserSession, error)
	UpdateSessionActivity(ctx context.Context, sessionToken string) error
	TerminateSession(ctx context.Context, sessionID int64, userID int64) error
	TerminateAllUserSessions(ctx context.Context, userID int64, exceptSessionID *int64) error
	CleanupExpiredSessions(ctx context.Context) error
	SetCurrentSession(ctx context.Context, sessionID int64, userID int64) error
	UpdateSessionTokens(ctx context.Context, sessionID int64, accessToken, refreshToken string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) SessionRepository {
	return &repository{
		db: db,
	}
}