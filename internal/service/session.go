package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"moveshare/internal/models"
	"moveshare/internal/repository/session"
	"strings"
	"time"
)

type SessionService interface {
	CreateSession(ctx context.Context, req *models.CreateSessionRequest, userID int64, accessToken, refreshToken string) (*models.UserSession, error)
	GetUserActiveSessions(ctx context.Context, userID int64) (*models.ActiveSessionsResponse, error)
	TerminateSession(ctx context.Context, sessionID int64, userID int64) error
	TerminateAllSessions(ctx context.Context, userID int64, exceptCurrentSession bool) error
	UpdateSessionActivity(ctx context.Context, sessionToken string) error
	UpdateSessionTokens(ctx context.Context, sessionID int64, accessToken, refreshToken string) error
}

type sessionService struct {
	sessionRepo session.SessionRepository
}

func NewSessionService(sessionRepo session.SessionRepository) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *sessionService) CreateSession(ctx context.Context, req *models.CreateSessionRequest, userID int64, accessToken, refreshToken string) (*models.UserSession, error) {
	deviceInfoJSON, err := json.Marshal(req.DeviceInfo)
	if err != nil {
		return nil, err
	}

	locationInfoJSON, err := json.Marshal(req.LocationInfo)
	if err != nil {
		return nil, err
	}

	session := &models.UserSession{
		UserID:       userID,
		SessionToken: accessToken, // Using access token as session identifier
		RefreshToken: &refreshToken,
		UserAgent:    &req.UserAgent,
		IPAddress:    &req.IPAddress,
		DeviceInfo:   deviceInfoJSON,
		LocationInfo: locationInfoJSON,
		IsCurrent:    true, // New session is current
	}

	// Set expiration time (e.g., 30 days from now)
	// Note: CreatedAt will be set by the database, so we use current time
	session.ExpiresAt = time.Now().AddDate(0, 0, 30)

	err = s.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	// Set this session as current (and unset others)
	err = s.sessionRepo.SetCurrentSession(ctx, session.ID, userID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionService) GetUserActiveSessions(ctx context.Context, userID int64) (*models.ActiveSessionsResponse, error) {
	log.Printf("Getting active sessions for user %d", userID)
	sessions, err := s.sessionRepo.GetUserSessions(ctx, userID)
	if err != nil {
		log.Printf("SessionService: Error getting sessions from repository for user %d: %v", userID, err)
		return nil, err
	}
	log.Printf("SessionService: Found %d sessions for user %d", len(sessions), userID)

	var responseSessions []models.UserSessionResponse
	
	// Initialize empty slice to avoid null in JSON
	if sessions == nil {
		sessions = []models.UserSession{}
	}

	for _, session := range sessions {
		// Parse device info
		var deviceInfo models.DeviceInfo
		if len(session.DeviceInfo) > 0 {
			json.Unmarshal(session.DeviceInfo, &deviceInfo)
		}

		// Parse location info
		var locationInfo models.LocationInfo
		if len(session.LocationInfo) > 0 {
			json.Unmarshal(session.LocationInfo, &locationInfo)
		}

		// Format device info
		deviceStr := formatDeviceInfo(deviceInfo)
		
		// Format location info
		locationStr := formatLocationInfo(locationInfo)

		ipAddress := ""
		if session.IPAddress != nil {
			ipAddress = *session.IPAddress
		}

		responseSession := models.UserSessionResponse{
			ID:           session.ID,
			DeviceInfo:   deviceStr,
			LocationInfo: locationStr,
			IPAddress:    ipAddress,
			IsCurrent:    session.IsCurrent,
			LastActivity: session.LastActivity,
			CreatedAt:    session.CreatedAt,
		}

		responseSessions = append(responseSessions, responseSession)
	}

	return &models.ActiveSessionsResponse{
		Sessions: responseSessions,
		Count:    len(responseSessions),
	}, nil
}

func (s *sessionService) TerminateSession(ctx context.Context, sessionID int64, userID int64) error {
	return s.sessionRepo.TerminateSession(ctx, sessionID, userID)
}

func (s *sessionService) TerminateAllSessions(ctx context.Context, userID int64, exceptCurrentSession bool) error {
	if exceptCurrentSession {
		// Find current session and keep it
		sessions, err := s.sessionRepo.GetUserSessions(ctx, userID)
		if err != nil {
			return err
		}

		var currentSessionID *int64
		for _, session := range sessions {
			if session.IsCurrent {
				currentSessionID = &session.ID
				break
			}
		}

		return s.sessionRepo.TerminateAllUserSessions(ctx, userID, currentSessionID)
	}

	return s.sessionRepo.TerminateAllUserSessions(ctx, userID, nil)
}

func (s *sessionService) UpdateSessionActivity(ctx context.Context, sessionToken string) error {
	return s.sessionRepo.UpdateSessionActivity(ctx, sessionToken)
}

func (s *sessionService) UpdateSessionTokens(ctx context.Context, sessionID int64, accessToken, refreshToken string) error {
	return s.sessionRepo.UpdateSessionTokens(ctx, sessionID, accessToken, refreshToken)
}

func formatDeviceInfo(device models.DeviceInfo) string {
	if device.Browser != "" && device.OS != "" {
		return fmt.Sprintf("%s on %s", device.Browser, device.OS)
	}
	if device.Browser != "" {
		return device.Browser
	}
	if device.OS != "" {
		return device.OS
	}
	return "Unknown Device"
}

func formatLocationInfo(location models.LocationInfo) string {
	parts := []string{}
	
	if location.City != "" {
		parts = append(parts, location.City)
	}
	if location.Region != "" {
		parts = append(parts, location.Region)
	}
	if location.Country != "" {
		parts = append(parts, location.Country)
	}

	if len(parts) > 0 {
		return strings.Join(parts, ", ")
	}

	return "Unknown Location"
}