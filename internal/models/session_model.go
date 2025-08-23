package models

import (
	"time"
	"encoding/json"
)

type UserSession struct {
	ID           int64           `json:"id" db:"id"`
	UserID       int64           `json:"user_id" db:"user_id"`
	SessionToken string          `json:"session_token" db:"session_token"`
	RefreshToken *string         `json:"refresh_token,omitempty" db:"refresh_token"`
	UserAgent    *string         `json:"user_agent" db:"user_agent"`
	IPAddress    *string         `json:"ip_address" db:"ip_address"`
	DeviceInfo   json.RawMessage `json:"device_info,omitempty" db:"device_info"`
	LocationInfo json.RawMessage `json:"location_info,omitempty" db:"location_info"`
	IsCurrent    bool            `json:"is_current" db:"is_current"`
	LastActivity time.Time       `json:"last_activity" db:"last_activity"`
	ExpiresAt    time.Time       `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

type DeviceInfo struct {
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	Device    string `json:"device"`
	Platform  string `json:"platform"`
}

type LocationInfo struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	IP      string `json:"ip"`
}

type ActiveSessionsResponse struct {
	Sessions []UserSessionResponse `json:"sessions"`
	Count    int                   `json:"count"`
}

type UserSessionResponse struct {
	ID           int64        `json:"id"`
	DeviceInfo   string       `json:"device_info"`    // Formatted as "Browser on OS"
	LocationInfo string       `json:"location_info"`  // Formatted as "City, Region"
	IPAddress    string       `json:"ip_address"`
	IsCurrent    bool         `json:"is_current"`
	LastActivity time.Time    `json:"last_activity"`
	CreatedAt    time.Time    `json:"created_at"`
}

type CreateSessionRequest struct {
	UserAgent    string       `json:"user_agent"`
	IPAddress    string       `json:"ip_address"`
	DeviceInfo   DeviceInfo   `json:"device_info"`
	LocationInfo LocationInfo `json:"location_info"`
}

type TerminateSessionRequest struct {
	SessionID int64 `json:"session_id"`
}