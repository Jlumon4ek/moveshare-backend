package models

import "time"

type EmailVerification struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Code      string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Used      bool      `json:"used"`
}

type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
}

type VerifyEmailCodeRequest struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
	Code  string `json:"code" binding:"required" example:"123456"`
}

type VerifyEmailCodeResponse struct {
	Message string `json:"message"`
	Valid   bool   `json:"valid"`
}