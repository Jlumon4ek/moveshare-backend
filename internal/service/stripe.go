package service

import (
	"context"
	"fmt"
	"log"
	"moveshare/internal/config"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/paymentmethod"
	"github.com/stripe/stripe-go/v82/setupintent"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeService interface {
	// Customer management
	CreateCustomer(ctx context.Context, userID int64, email, name string) (string, error)
	GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error)

	// Payment Methods
	CreateSetupIntent(ctx context.Context, customerID string) (*stripe.SetupIntent, error)
	AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) (*stripe.PaymentMethod, error)
	DetachPaymentMethod(ctx context.Context, paymentMethodID string) (*stripe.PaymentMethod, error)
	ListPaymentMethods(ctx context.Context, customerID string) ([]*stripe.PaymentMethod, error)
	SetDefaultPaymentMethod(ctx context.Context, customerID, paymentMethodID string) error

	// Payments
	CreatePaymentIntent(ctx context.Context, amount int64, currency, customerID, paymentMethodID, description string) (*stripe.PaymentIntent, error)
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error)
	GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error)
	CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error)

	// Webhook
	ConstructEvent(payload []byte, header string) (stripe.Event, error)
}

type stripeService struct {
	config *config.StripeConfig
}

func NewStripeService(cfg *config.StripeConfig) StripeService {
	// Устанавливаем API ключ Stripe
	stripe.Key = cfg.PrivateKey

	log.Printf("Stripe service initialized with key: %s...", cfg.PrivateKey[:20])

	return &stripeService{
		config: cfg,
	}
}

// Customer management
func (s *stripeService) CreateCustomer(ctx context.Context, userID int64, email, name string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	c, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	return c.ID, nil
}

func (s *stripeService) GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	c, err := customer.Get(customerID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Stripe customer: %w", err)
	}

	return c, nil
}

// Payment Methods
func (s *stripeService) CreateSetupIntent(ctx context.Context, customerID string) (*stripe.SetupIntent, error) {
	params := &stripe.SetupIntentParams{
		Customer: stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Usage: stripe.String("off_session"),
	}

	si, err := setupintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create setup intent: %w", err)
	}

	return si, nil
}

func (s *stripeService) AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}

	pm, err := paymentmethod.Attach(paymentMethodID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method: %w", err)
	}

	return pm, nil
}

func (s *stripeService) DetachPaymentMethod(ctx context.Context, paymentMethodID string) (*stripe.PaymentMethod, error) {
	pm, err := paymentmethod.Detach(paymentMethodID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to detach payment method: %w", err)
	}

	return pm, nil
}

func (s *stripeService) ListPaymentMethods(ctx context.Context, customerID string) ([]*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}

	iter := paymentmethod.List(params)
	var paymentMethods []*stripe.PaymentMethod

	for iter.Next() {
		paymentMethods = append(paymentMethods, iter.PaymentMethod())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}

	return paymentMethods, nil
}

func (s *stripeService) SetDefaultPaymentMethod(ctx context.Context, customerID, paymentMethodID string) error {
	params := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethodID),
		},
	}

	_, err := customer.Update(customerID, params)
	if err != nil {
		return fmt.Errorf("failed to set default payment method: %w", err)
	}

	return nil
}

// Payments
func (s *stripeService) CreatePaymentIntent(ctx context.Context, amount int64, currency, customerID, paymentMethodID, description string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		PaymentMethod:      stripe.String(paymentMethodID),
		Description:        stripe.String(description),
		ConfirmationMethod: stripe.String("manual"),
		Confirm:            stripe.Bool(false), // ✅ ИЗМЕНЕНИЕ: не подтверждаем автоматически
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi, nil
}

func (s *stripeService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{}

	pi, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return pi, nil
}

func (s *stripeService) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return pi, nil
}

func (s *stripeService) CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Cancel(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment intent: %w", err)
	}

	return pi, nil
}

// Webhook
func (s *stripeService) ConstructEvent(payload []byte, header string) (stripe.Event, error) {
	// ✅ ИСПРАВЛЕНИЕ: Используем webhook.ConstructEvent для v82
	event, err := webhook.ConstructEvent(payload, header, s.config.WebhookSecret)
	if err != nil {
		return event, fmt.Errorf("failed to construct stripe event: %w", err)
	}

	return event, nil
}
