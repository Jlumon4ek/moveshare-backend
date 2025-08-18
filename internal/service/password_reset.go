package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"moveshare/internal/models"
	"moveshare/internal/repository/password_reset"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordResetService interface {
	SendResetCode(ctx context.Context, email string) error
	VerifyResetCode(ctx context.Context, email, code string) error
	ResetPassword(ctx context.Context, email, code, newPassword string) error
}

type passwordResetService struct {
	passwordResetRepo password_reset.PasswordResetRepository
	emailService      EmailService
}

func NewPasswordResetService(passwordResetRepo password_reset.PasswordResetRepository, emailService EmailService) PasswordResetService {
	return &passwordResetService{
		passwordResetRepo: passwordResetRepo,
		emailService:      emailService,
	}
}

func (s *passwordResetService) SendResetCode(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.passwordResetRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Generate 6-digit code
	code, err := generateResetCode()
	if err != nil {
		return fmt.Errorf("failed to generate reset code: %w", err)
	}

	// Create reset code record
	resetCode := &models.PasswordResetCode{
		UserID:    int(user.ID),
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute), // Code expires in 15 minutes
		Used:      false,
	}

	err = s.passwordResetRepo.CreateResetCode(ctx, resetCode)
	if err != nil {
		return fmt.Errorf("failed to save reset code: %w", err)
	}

	// Send email
	err = s.emailService.SendPasswordResetCode(email, code)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *passwordResetService) VerifyResetCode(ctx context.Context, email, code string) error {
	// Find valid reset code
	_, err := s.passwordResetRepo.GetValidResetCode(ctx, email, code)
	if err != nil {
		return fmt.Errorf("invalid or expired reset code")
	}

	return nil
}

func (s *passwordResetService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	// Find valid reset code
	resetCode, err := s.passwordResetRepo.GetValidResetCode(ctx, email, code)
	if err != nil {
		return fmt.Errorf("invalid or expired reset code")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	err = s.passwordResetRepo.UpdateUserPassword(ctx, resetCode.UserID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark code as used
	err = s.passwordResetRepo.MarkCodeAsUsed(ctx, resetCode.ID)
	if err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	return nil
}

func generateResetCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += digit.String()
	}
	return code, nil
}