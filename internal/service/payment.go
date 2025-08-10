// internal/service/payment.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"moveshare/internal/models"
	"moveshare/internal/repository/payment"
	"strings"

	"github.com/stripe/stripe-go/v82"
)

type PaymentService interface {
	// Payment Methods
	AddPaymentMethod(ctx context.Context, userID int64, paymentMethodID string) (*models.UserPaymentMethod, error)
	GetUserPaymentMethods(ctx context.Context, userID int64) ([]models.UserPaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, userID, paymentMethodID int64) error
	SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID int64) error
	GetDefaultPaymentMethod(ctx context.Context, userID int64) (*models.UserPaymentMethod, error)

	// Customer Management
	EnsureStripeCustomer(ctx context.Context, userID int64) (string, error)
	GetOrCreateStripeCustomer(ctx context.Context, userID int64, email, name string) (string, error)

	// Payments
	CreatePayment(ctx context.Context, userID int64, req *models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
	ConfirmPayment(ctx context.Context, paymentIntentID string) (*models.ConfirmPaymentResponse, error)
	GetUserPayments(ctx context.Context, userID int64, limit, offset int) ([]models.Payment, error)

	// Webhook
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
}

type paymentService struct {
	paymentRepo   payment.PaymentRepository
	stripeService StripeService
	userService   UserService
}

func NewPaymentService(
	paymentRepo payment.PaymentRepository,
	stripeService StripeService,
	userService UserService,
) PaymentService {
	return &paymentService{
		paymentRepo:   paymentRepo,
		stripeService: stripeService,
		userService:   userService,
	}
}

// Customer Management
func (s *paymentService) EnsureStripeCustomer(ctx context.Context, userID int64) (string, error) {
	// Проверяем есть ли уже Stripe customer
	customerID, err := s.paymentRepo.GetUserStripeCustomerID(ctx, userID)
	if err != nil {
		return "", err
	}

	if customerID != "" {
		return customerID, nil
	}

	// Получаем информацию о пользователе
	user, err := s.userService.FindUserByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	// Создаем Stripe customer
	customerID, err = s.stripeService.CreateCustomer(ctx, userID, user.Email, user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Сохраняем customer ID в БД
	err = s.paymentRepo.UpdateUserStripeCustomerID(ctx, userID, customerID)
	if err != nil {
		return "", fmt.Errorf("failed to save customer ID: %w", err)
	}

	return customerID, nil
}

func (s *paymentService) GetOrCreateStripeCustomer(ctx context.Context, userID int64, email, name string) (string, error) {
	return s.EnsureStripeCustomer(ctx, userID)
}

// Payment Methods
func (s *paymentService) AddPaymentMethod(ctx context.Context, userID int64, paymentMethodID string) (*models.UserPaymentMethod, error) {
	// Убеждаемся что у пользователя есть Stripe customer
	customerID, err := s.EnsureStripeCustomer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure Stripe customer: %w", err)
	}

	// Проверяем что payment method не добавлен уже
	existingPM, err := s.paymentRepo.GetPaymentMethodByStripeID(ctx, userID, paymentMethodID)
	if err == nil && existingPM != nil {
		return nil, fmt.Errorf("payment method already added")
	}

	// Привязываем payment method к customer в Stripe
	stripePaymentMethod, err := s.stripeService.AttachPaymentMethod(ctx, paymentMethodID, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method to customer: %w", err)
	}

	// Проверяем что это карта
	if stripePaymentMethod.Type != "card" {
		return nil, fmt.Errorf("unsupported payment method type: %s", stripePaymentMethod.Type)
	}

	// Создаем запись в БД
	userPaymentMethod := &models.UserPaymentMethod{
		UserID:                userID,
		StripePaymentMethodID: stripePaymentMethod.ID,
		StripeCustomerID:      customerID,
		CardLast4:             stripePaymentMethod.Card.Last4,
		CardBrand:             strings.Title(string(stripePaymentMethod.Card.Brand)),
		CardExpMonth:          int(stripePaymentMethod.Card.ExpMonth),
		CardExpYear:           int(stripePaymentMethod.Card.ExpYear),
		IsDefault:             false, // Будет установлено в репозитории если это первая карта
		IsActive:              true,
	}

	err = s.paymentRepo.SavePaymentMethod(ctx, userPaymentMethod)
	if err != nil {
		// Откатываем изменения в Stripe если не удалось сохранить в БД
		_, detachErr := s.stripeService.DetachPaymentMethod(ctx, paymentMethodID)
		if detachErr != nil {
			// Логируем ошибку, но возвращаем оригинальную
			fmt.Printf("Failed to detach payment method after DB error: %v\n", detachErr)
		}
		return nil, fmt.Errorf("failed to save payment method: %w", err)
	}

	return userPaymentMethod, nil
}

func (s *paymentService) GetUserPaymentMethods(ctx context.Context, userID int64) ([]models.UserPaymentMethod, error) {
	return s.paymentRepo.GetUserPaymentMethods(ctx, userID)
}

func (s *paymentService) DeletePaymentMethod(ctx context.Context, userID, paymentMethodID int64) error {
	// Получаем информацию о payment method
	paymentMethod, err := s.paymentRepo.GetPaymentMethodByID(ctx, userID, paymentMethodID)
	if err != nil {
		return fmt.Errorf("payment method not found: %w", err)
	}

	// Отвязываем от Stripe customer
	_, err = s.stripeService.DetachPaymentMethod(ctx, paymentMethod.StripePaymentMethodID)
	if err != nil {
		// Логируем ошибку, но продолжаем удаление из БД
		fmt.Printf("Failed to detach payment method from Stripe: %v\n", err)
	}

	// Удаляем из БД
	err = s.paymentRepo.DeletePaymentMethod(ctx, userID, paymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	return nil
}

func (s *paymentService) SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID int64) error {
	// Получаем информацию о payment method
	paymentMethod, err := s.paymentRepo.GetPaymentMethodByID(ctx, userID, paymentMethodID)
	if err != nil {
		return fmt.Errorf("payment method not found: %w", err)
	}

	// Устанавливаем default в Stripe
	err = s.stripeService.SetDefaultPaymentMethod(ctx, paymentMethod.StripeCustomerID, paymentMethod.StripePaymentMethodID)
	if err != nil {
		// Логируем ошибку, но продолжаем обновление в БД
		fmt.Printf("Failed to set default payment method in Stripe: %v\n", err)
	}

	// Устанавливаем default в БД
	err = s.paymentRepo.SetDefaultPaymentMethod(ctx, userID, paymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to set default payment method: %w", err)
	}

	return nil
}

func (s *paymentService) GetDefaultPaymentMethod(ctx context.Context, userID int64) (*models.UserPaymentMethod, error) {
	return s.paymentRepo.GetDefaultPaymentMethod(ctx, userID)
}

// Payments
func (s *paymentService) CreatePayment(ctx context.Context, userID int64, req *models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	// Убеждаемся что у пользователя есть Stripe customer
	customerID, err := s.EnsureStripeCustomer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure Stripe customer: %w", err)
	}

	// Определяем payment method
	var paymentMethod *models.UserPaymentMethod
	if req.PaymentMethodID != nil {
		// Используем указанный payment method
		paymentMethod, err = s.paymentRepo.GetPaymentMethodByID(ctx, userID, *req.PaymentMethodID)
		if err != nil {
			return nil, fmt.Errorf("payment method not found: %w", err)
		}
	} else {
		// Используем default payment method
		paymentMethod, err = s.paymentRepo.GetDefaultPaymentMethod(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("no default payment method found: %w", err)
		}
	}

	// Создаем Payment Intent в Stripe
	paymentIntent, err := s.stripeService.CreatePaymentIntent(
		ctx,
		req.AmountCents,
		"usd", // TODO: сделать конфигурируемым
		customerID,
		paymentMethod.StripePaymentMethodID,
		req.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Сохраняем платеж в БД
	payment := &models.Payment{
		UserID:                userID,
		JobID:                 req.JobID,
		StripePaymentIntentID: paymentIntent.ID,
		StripePaymentMethodID: paymentMethod.StripePaymentMethodID,
		StripeCustomerID:      customerID,
		AmountCents:           req.AmountCents,
		Currency:              "usd",
		Status:                string(paymentIntent.Status),
		Description:           req.Description,
	}

	err = s.paymentRepo.SavePayment(ctx, payment)
	if err != nil {
		// Отменяем payment intent если не удалось сохранить в БД
		_, cancelErr := s.stripeService.CancelPaymentIntent(ctx, paymentIntent.ID)
		if cancelErr != nil {
			fmt.Printf("Failed to cancel payment intent after DB error: %v\n", cancelErr)
		}
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return &models.CreatePaymentResponse{
		PaymentIntentID:      paymentIntent.ID,
		ClientSecret:         paymentIntent.ClientSecret,
		Status:               string(paymentIntent.Status),
		RequiresConfirmation: paymentIntent.Status == stripe.PaymentIntentStatusRequiresConfirmation,
		Success:              true,
	}, nil
}

func (s *paymentService) ConfirmPayment(ctx context.Context, paymentIntentID string) (*models.ConfirmPaymentResponse, error) {
	// Подтверждаем платеж в Stripe
	paymentIntent, err := s.stripeService.ConfirmPaymentIntent(ctx, paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	// Получаем платеж из БД
	payment, err := s.paymentRepo.GetPaymentByStripeIntentID(ctx, paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	var failureReason string
	if paymentIntent.LastPaymentError != nil && paymentIntent.LastPaymentError.Err != nil {
		failureReason = paymentIntent.LastPaymentError.Err.Error()
	}

	err = s.paymentRepo.UpdatePaymentStatus(ctx, payment.ID, string(paymentIntent.Status), failureReason)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	return &models.ConfirmPaymentResponse{
		PaymentID: payment.ID,
		Status:    string(paymentIntent.Status),
		Success:   paymentIntent.Status == stripe.PaymentIntentStatusSucceeded,
		Message:   getPaymentStatusMessage(paymentIntent.Status),
	}, nil
}

func (s *paymentService) GetUserPayments(ctx context.Context, userID int64, limit, offset int) ([]models.Payment, error) {
	return s.paymentRepo.GetUserPayments(ctx, userID, limit, offset)
}

// Webhook
func (s *paymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := s.stripeService.ConstructEvent(payload, signature)
	if err != nil {
		return fmt.Errorf("failed to construct webhook event: %w", err)
	}

	switch event.Type {
	case "payment_intent.succeeded":
		return s.handlePaymentIntentSucceeded(ctx, event)
	case "payment_intent.payment_failed":
		return s.handlePaymentIntentFailed(ctx, event)
	case "payment_method.attached":
		return s.handlePaymentMethodAttached(ctx, event)
	case "payment_method.detached":
		return s.handlePaymentMethodDetached(ctx, event)
	default:
		// Игнорируем неизвестные события
		fmt.Printf("Received unhandled webhook event: %s\n", event.Type)
		return nil
	}
}

func (s *paymentService) handlePaymentIntentSucceeded(ctx context.Context, event stripe.Event) error {
	// ✅ ИСПРАВЛЕНИЕ: Правильный способ получения данных из webhook для v82
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("failed to parse payment intent from webhook: %w", err)
	}

	// Получаем полную информацию о payment intent
	fullPaymentIntent, err := s.stripeService.GetPaymentIntent(ctx, paymentIntent.ID)
	if err != nil {
		return fmt.Errorf("failed to get payment intent: %w", err)
	}

	// Обновляем статус в БД
	payment, err := s.paymentRepo.GetPaymentByStripeIntentID(ctx, fullPaymentIntent.ID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	err = s.paymentRepo.UpdatePaymentStatus(ctx, payment.ID, string(fullPaymentIntent.Status), "")
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// TODO: Активировать job если это платеж за размещение
	// if payment.JobID > 0 {
	//     err = s.jobService.ActivateJob(ctx, payment.JobID)
	//     if err != nil {
	//         log.Printf("Failed to activate job %d: %v", payment.JobID, err)
	//     }
	// }

	return nil
}

func (s *paymentService) handlePaymentIntentFailed(ctx context.Context, event stripe.Event) error {
	// ✅ ИСПРАВЛЕНИЕ: Правильный способ получения данных из webhook для v82
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("failed to parse payment intent from webhook: %w", err)
	}

	fullPaymentIntent, err := s.stripeService.GetPaymentIntent(ctx, paymentIntent.ID)
	if err != nil {
		return fmt.Errorf("failed to get payment intent: %w", err)
	}

	payment, err := s.paymentRepo.GetPaymentByStripeIntentID(ctx, fullPaymentIntent.ID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	var failureReason string
	if fullPaymentIntent.LastPaymentError != nil && fullPaymentIntent.LastPaymentError.Err != nil {
		failureReason = fullPaymentIntent.LastPaymentError.Err.Error()
	}

	err = s.paymentRepo.UpdatePaymentStatus(ctx, payment.ID, string(fullPaymentIntent.Status), failureReason)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

func (s *paymentService) handlePaymentMethodAttached(ctx context.Context, event stripe.Event) error {
	// Этот webhook можно использовать для дополнительной валидации
	// или логирования, но основная логика уже в AddPaymentMethod
	return nil
}

func (s *paymentService) handlePaymentMethodDetached(ctx context.Context, event stripe.Event) error {
	// Этот webhook можно использовать для дополнительной валидации
	// или логирования, но основная логика уже в DeletePaymentMethod
	return nil
}

func getPaymentStatusMessage(status stripe.PaymentIntentStatus) string {
	switch status {
	case stripe.PaymentIntentStatusSucceeded:
		return "Payment successful"
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		return "Payment requires a valid payment method"
	case stripe.PaymentIntentStatusRequiresConfirmation:
		return "Payment requires confirmation"
	case stripe.PaymentIntentStatusRequiresAction:
		return "Payment requires additional action"
	case stripe.PaymentIntentStatusProcessing:
		return "Payment is processing"
	case stripe.PaymentIntentStatusCanceled:
		return "Payment was canceled"
	default:
		return "Unknown payment status"
	}
}
