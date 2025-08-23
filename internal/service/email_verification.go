package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"moveshare/internal/models"
	"moveshare/internal/repository/email_verification"
	"time"
)

type EmailVerificationService interface {
	SendVerificationCode(ctx context.Context, email string) error
	VerifyEmailCode(ctx context.Context, email, code string) error
}

type emailVerificationService struct {
	emailVerificationRepo email_verification.EmailVerificationRepository
	emailService          EmailService
}

func NewEmailVerificationService(emailVerificationRepo email_verification.EmailVerificationRepository, emailService EmailService) EmailVerificationService {
	return &emailVerificationService{
		emailVerificationRepo: emailVerificationRepo,
		emailService:          emailService,
	}
}

func (s *emailVerificationService) SendVerificationCode(ctx context.Context, email string) error {
	// Generate 6-digit verification code
	code, err := s.generateVerificationCode()
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Create verification record
	verification := &models.EmailVerification{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute), // Code expires in 15 minutes
		CreatedAt: time.Now(),
		Used:      false,
	}

	err = s.emailVerificationRepo.CreateVerificationCode(ctx, verification)
	if err != nil {
		return fmt.Errorf("failed to create verification code: %w", err)
	}

	// Send email
	subject := "✨ Email Verification Code - MoveShare"
	
	// Beautiful HTML version
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification - MoveShare</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);">
        
        <!-- Header -->
        <div style="background: linear-gradient(135deg, #10B981 0%%, #059669 100%%); padding: 40px 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700; text-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                🚚 MoveShare
            </h1>
            <p style="color: #A7F3D0; margin: 8px 0 0 0; font-size: 16px; opacity: 0.9;">
                Welcome! Please verify your email
            </p>
        </div>

        <!-- Content -->
        <div style="padding: 40px 30px;">
            <h2 style="color: #1F2937; margin: 0 0 24px 0; font-size: 24px; font-weight: 600;">
                Email Verification Required
            </h2>
            
            <p style="color: #4B5563; margin: 0 0 24px 0; font-size: 16px; line-height: 1.6;">
                Welcome to MoveShare! To complete your registration, please verify your email address using the code below.
            </p>

            <!-- Code Box -->
            <div style="background: linear-gradient(135deg, #10B98110 0%%, #05966920 100%%); border: 2px solid #10B981; border-radius: 12px; padding: 24px; margin: 32px 0; text-align: center;">
                <p style="color: #374151; margin: 0 0 12px 0; font-size: 14px; font-weight: 600; text-transform: uppercase; letter-spacing: 1px;">
                    Your Verification Code
                </p>
                <div style="background: #ffffff; border-radius: 8px; padding: 16px; margin: 12px 0; border: 1px solid #E5E7EB;">
                    <span style="font-family: 'Courier New', monospace; font-size: 32px; font-weight: 700; color: #10B981; letter-spacing: 4px;">
                        %s
                    </span>
                </div>
                <p style="color: #6B7280; margin: 12px 0 0 0; font-size: 14px;">
                    ⏰ This code expires in <strong>15 minutes</strong>
                </p>
            </div>

            <div style="background: #DBEAFE; border-left: 4px solid #3B82F6; padding: 16px; margin: 24px 0; border-radius: 0 8px 8px 0;">
                <p style="color: #1E40AF; margin: 0; font-size: 14px;">
                    <strong>Next Steps:</strong> Enter this code in the registration form to complete your account setup and start using MoveShare!
                </p>
            </div>

            <div style="background: #FEF3C7; border-left: 4px solid #F59E0B; padding: 16px; margin: 24px 0; border-radius: 0 8px 8px 0;">
                <p style="color: #92400E; margin: 0; font-size: 14px;">
                    <strong>Security Notice:</strong> If you didn't sign up for MoveShare, please ignore this email. No account will be created.
                </p>
            </div>

            <p style="color: #4B5563; margin: 24px 0 0 0; font-size: 16px; line-height: 1.6;">
                Thanks for joining us!<br>
                <strong style="color: #10B981;">The MoveShare Team</strong>
            </p>
        </div>

        <!-- Footer -->
        <div style="background: #F9FAFB; padding: 24px 30px; border-top: 1px solid #E5E7EB; text-align: center;">
            <p style="color: #6B7280; margin: 0; font-size: 12px;">
                © 2024 MoveShare. All rights reserved.<br>
                This is an automated message, please do not reply.
            </p>
        </div>
    </div>
</body>
</html>`, code)

	err = s.emailService.SendEmail(email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (s *emailVerificationService) VerifyEmailCode(ctx context.Context, email, code string) error {
	verification, err := s.emailVerificationRepo.GetVerificationByEmailAndCode(ctx, email, code)
	if err != nil {
		return fmt.Errorf("invalid or expired verification code")
	}

	// Mark code as used
	err = s.emailVerificationRepo.MarkCodeAsUsed(ctx, verification.ID)
	if err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	return nil
}

func (s *emailVerificationService) generateVerificationCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)
	
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}
	
	return string(code), nil
}