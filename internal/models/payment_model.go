// internal/models/payment_model.go
package models

import "time"

// UserPaymentMethod представляет сохраненную карту пользователя
type UserPaymentMethod struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	StripePaymentMethodID string    `json:"stripe_payment_method_id"`
	StripeCustomerID      string    `json:"stripe_customer_id,omitempty"`
	CardLast4             string    `json:"card_last4"`
	CardBrand             string    `json:"card_brand"`
	CardExpMonth          int       `json:"card_exp_month"`
	CardExpYear           int       `json:"card_exp_year"`
	IsDefault             bool      `json:"is_default"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Payment представляет платеж
type Payment struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	JobID                 int64     `json:"job_id,omitempty"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
	StripePaymentMethodID string    `json:"stripe_payment_method_id"`
	StripeCustomerID      string    `json:"stripe_customer_id"`
	AmountCents           int64     `json:"amount_cents"`
	Currency              string    `json:"currency"`
	Status                string    `json:"status"`
	Description           string    `json:"description,omitempty"`
	FailureReason         string    `json:"failure_reason,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// User расширение для Stripe Customer ID
type UserWithStripe struct {
	User
	StripeCustomerID string `json:"stripe_customer_id,omitempty"`
}

// DTO для API ответов
type PaymentMethodResponse struct {
	ID           int64     `json:"id"`
	CardLast4    string    `json:"card_last4"`
	CardBrand    string    `json:"card_brand"`
	CardExpMonth int       `json:"card_exp_month"`
	CardExpYear  int       `json:"card_exp_year"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
}

type AddCardRequest struct {
	PaymentMethodID string `json:"payment_method_id" binding:"required" example:"pm_1234567890"`
}

type AddCardResponse struct {
	PaymentMethod PaymentMethodResponse `json:"payment_method"`
	Message       string                `json:"message"`
	Success       bool                  `json:"success"`
}

type CreatePaymentRequest struct {
	JobID           int64  `json:"job_id" binding:"required" example:"123"`
	PaymentMethodID *int64 `json:"payment_method_id,omitempty" example:"456"`      // Опционально, если не указано - берем default
	AmountCents     int64  `json:"amount_cents" binding:"required" example:"2999"` // $29.99
	Description     string `json:"description,omitempty" example:"Payment for job posting"`
}

type CreatePaymentResponse struct {
	PaymentIntentID      string `json:"payment_intent_id"`
	ClientSecret         string `json:"client_secret"`
	Status               string `json:"status"`
	RequiresConfirmation bool   `json:"requires_confirmation"`
	Success              bool   `json:"success"`
}

type ConfirmPaymentRequest struct {
	PaymentIntentID string `json:"payment_intent_id" binding:"required" example:"pi_1234567890"`
}

type ConfirmPaymentResponse struct {
	PaymentID int64  `json:"payment_id"`
	Status    string `json:"status"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}
