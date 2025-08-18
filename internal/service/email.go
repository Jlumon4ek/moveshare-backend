package service

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

type EmailService interface {
	SendPasswordResetCode(to, code string) error
}

type emailService struct {
	client *resend.Client
	from   string
}

func NewEmailService() EmailService {
	apiKey := getEnv("RESEND_API_KEY", "re_2Bae7Xuc_QAEoCc2xUHBz1BDqcwU8UhTY")
	client := resend.NewClient(apiKey)
	
	return &emailService{
		client: client,
		from:   getEnv("EMAIL_FROM", "admin@themoveshare.com"),
	}
}

func (e *emailService) SendPasswordResetCode(to, code string) error {
	subject := "🔐 Password Reset Code - MoveShare"
	
	// Plain text version for fallback
	text := fmt.Sprintf(`Hello,

You have requested to reset your password for your MoveShare account.

Your password reset code is: %s

This code will expire in 15 minutes.

If you did not request this password reset, please ignore this email.

Best regards,
The MoveShare Team`, code)

	// Beautiful HTML version
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset - MoveShare</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);">
        
        <!-- Header -->
        <div style="background: linear-gradient(135deg, #60A5FA 0%%, #3B82F6 100%%); padding: 40px 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700; text-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                🚚 MoveShare
            </h1>
            <p style="color: #E0F2FE; margin: 8px 0 0 0; font-size: 16px; opacity: 0.9;">
                Secure Password Reset
            </p>
        </div>

        <!-- Content -->
        <div style="padding: 40px 30px;">
            <h2 style="color: #1F2937; margin: 0 0 24px 0; font-size: 24px; font-weight: 600;">
                Password Reset Request
            </h2>
            
            <p style="color: #4B5563; margin: 0 0 24px 0; font-size: 16px; line-height: 1.6;">
                Hello! We received a request to reset your password for your MoveShare account. 
            </p>

            <!-- Code Box -->
            <div style="background: linear-gradient(135deg, #60A5FA10 0%%, #3B82F620 100%%); border: 2px solid #60A5FA; border-radius: 12px; padding: 24px; margin: 32px 0; text-align: center;">
                <p style="color: #374151; margin: 0 0 12px 0; font-size: 14px; font-weight: 600; text-transform: uppercase; letter-spacing: 1px;">
                    Your Reset Code
                </p>
                <div style="background: #ffffff; border-radius: 8px; padding: 16px; margin: 12px 0; border: 1px solid #E5E7EB;">
                    <span style="font-family: 'Courier New', monospace; font-size: 32px; font-weight: 700; color: #60A5FA; letter-spacing: 4px;">
                        %s
                    </span>
                </div>
                <p style="color: #6B7280; margin: 12px 0 0 0; font-size: 14px;">
                    ⏰ This code expires in <strong>15 minutes</strong>
                </p>
            </div>

            <div style="background: #FEF3C7; border-left: 4px solid #F59E0B; padding: 16px; margin: 24px 0; border-radius: 0 8px 8px 0;">
                <p style="color: #92400E; margin: 0; font-size: 14px;">
                    <strong>Security Notice:</strong> If you didn't request this password reset, please ignore this email. Your account remains secure.
                </p>
            </div>

            <p style="color: #4B5563; margin: 24px 0 0 0; font-size: 16px; line-height: 1.6;">
                Best regards,<br>
                <strong style="color: #60A5FA;">The MoveShare Team</strong>
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

	params := &resend.SendEmailRequest{
		From:    e.from,
		To:      []string{to},
		Subject: subject,
		Text:    text,
		Html:    html,
	}

	_, err := e.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
