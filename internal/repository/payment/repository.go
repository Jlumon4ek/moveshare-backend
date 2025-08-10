// internal/repository/payment/repository.go
package payment

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	// User Customer Management
	GetUserStripeCustomerID(ctx context.Context, userID int64) (string, error)
	UpdateUserStripeCustomerID(ctx context.Context, userID int64, customerID string) error

	// Payment Methods
	SavePaymentMethod(ctx context.Context, paymentMethod *models.UserPaymentMethod) error
	GetUserPaymentMethods(ctx context.Context, userID int64) ([]models.UserPaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, userID, paymentMethodID int64) (*models.UserPaymentMethod, error)
	GetPaymentMethodByStripeID(ctx context.Context, userID int64, stripePaymentMethodID string) (*models.UserPaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, userID, paymentMethodID int64) error
	SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID int64) error
	GetDefaultPaymentMethod(ctx context.Context, userID int64) (*models.UserPaymentMethod, error)

	SavePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByStripeIntentID(ctx context.Context, stripePaymentIntentID string) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int64, status, failureReason string) error
	GetUserPayments(ctx context.Context, userID int64, limit, offset int) ([]models.Payment, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) PaymentRepository {
	return &repository{db: db}
}

// internal/repository/payment/customer.go
func (r *repository) GetUserStripeCustomerID(ctx context.Context, userID int64) (string, error) {
	query := `SELECT stripe_customer_id FROM users WHERE id = $1`

	var customerID *string
	err := r.db.QueryRow(ctx, query, userID).Scan(&customerID)
	if err != nil {
		return "", err
	}

	if customerID == nil {
		return "", nil
	}

	return *customerID, nil
}

func (r *repository) UpdateUserStripeCustomerID(ctx context.Context, userID int64, customerID string) error {
	query := `UPDATE users SET stripe_customer_id = $1 WHERE id = $2`

	_, err := r.db.Exec(ctx, query, customerID, userID)
	return err
}

// internal/repository/payment/payment_methods.go
func (r *repository) SavePaymentMethod(ctx context.Context, paymentMethod *models.UserPaymentMethod) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Если это будет default карта, сначала убираем default у других
	if paymentMethod.IsDefault {
		updateQuery := `
			UPDATE user_payment_methods 
			SET is_default = false, updated_at = NOW()
			WHERE user_id = $1 AND is_active = true
		`
		_, err = tx.Exec(ctx, updateQuery, paymentMethod.UserID)
		if err != nil {
			return err
		}
	}

	// Если у пользователя нет других карт, делаем эту default
	if !paymentMethod.IsDefault {
		countQuery := `
			SELECT COUNT(*) 
			FROM user_payment_methods 
			WHERE user_id = $1 AND is_active = true
		`
		var count int
		err = tx.QueryRow(ctx, countQuery, paymentMethod.UserID).Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			paymentMethod.IsDefault = true
		}
	}

	// Вставляем новую карту
	insertQuery := `
		INSERT INTO user_payment_methods (
			user_id, stripe_payment_method_id, stripe_customer_id, 
			card_last4, card_brand, card_exp_month, card_exp_year, 
			is_default, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, insertQuery,
		paymentMethod.UserID,
		paymentMethod.StripePaymentMethodID,
		paymentMethod.StripeCustomerID,
		paymentMethod.CardLast4,
		paymentMethod.CardBrand,
		paymentMethod.CardExpMonth,
		paymentMethod.CardExpYear,
		paymentMethod.IsDefault,
		paymentMethod.IsActive,
	).Scan(&paymentMethod.ID, &paymentMethod.CreatedAt, &paymentMethod.UpdatedAt)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *repository) GetUserPaymentMethods(ctx context.Context, userID int64) ([]models.UserPaymentMethod, error) {
	query := `
		SELECT id, user_id, stripe_payment_method_id, stripe_customer_id,
		       card_last4, card_brand, card_exp_month, card_exp_year,
		       is_default, is_active, created_at, updated_at
		FROM user_payment_methods
		WHERE user_id = $1 AND is_active = true
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentMethods []models.UserPaymentMethod
	for rows.Next() {
		var pm models.UserPaymentMethod
		err := rows.Scan(
			&pm.ID, &pm.UserID, &pm.StripePaymentMethodID, &pm.StripeCustomerID,
			&pm.CardLast4, &pm.CardBrand, &pm.CardExpMonth, &pm.CardExpYear,
			&pm.IsDefault, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		paymentMethods = append(paymentMethods, pm)
	}

	return paymentMethods, rows.Err()
}

func (r *repository) GetPaymentMethodByID(ctx context.Context, userID, paymentMethodID int64) (*models.UserPaymentMethod, error) {
	query := `
		SELECT id, user_id, stripe_payment_method_id, stripe_customer_id,
		       card_last4, card_brand, card_exp_month, card_exp_year,
		       is_default, is_active, created_at, updated_at
		FROM user_payment_methods
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`

	var pm models.UserPaymentMethod
	err := r.db.QueryRow(ctx, query, paymentMethodID, userID).Scan(
		&pm.ID, &pm.UserID, &pm.StripePaymentMethodID, &pm.StripeCustomerID,
		&pm.CardLast4, &pm.CardBrand, &pm.CardExpMonth, &pm.CardExpYear,
		&pm.IsDefault, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &pm, nil
}

func (r *repository) GetPaymentMethodByStripeID(ctx context.Context, userID int64, stripePaymentMethodID string) (*models.UserPaymentMethod, error) {
	query := `
		SELECT id, user_id, stripe_payment_method_id, stripe_customer_id,
		       card_last4, card_brand, card_exp_month, card_exp_year,
		       is_default, is_active, created_at, updated_at
		FROM user_payment_methods
		WHERE stripe_payment_method_id = $1 AND user_id = $2 AND is_active = true
	`

	var pm models.UserPaymentMethod
	err := r.db.QueryRow(ctx, query, stripePaymentMethodID, userID).Scan(
		&pm.ID, &pm.UserID, &pm.StripePaymentMethodID, &pm.StripeCustomerID,
		&pm.CardLast4, &pm.CardBrand, &pm.CardExpMonth, &pm.CardExpYear,
		&pm.IsDefault, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &pm, nil
}

func (r *repository) DeletePaymentMethod(ctx context.Context, userID, paymentMethodID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Проверяем была ли карта default
	var wasDefault bool
	checkQuery := `
		SELECT is_default 
		FROM user_payment_methods 
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`
	err = tx.QueryRow(ctx, checkQuery, paymentMethodID, userID).Scan(&wasDefault)
	if err != nil {
		return err
	}

	// Деактивируем карту
	deleteQuery := `
		UPDATE user_payment_methods 
		SET is_active = false, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	result, err := tx.Exec(ctx, deleteQuery, paymentMethodID, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return err // Not found or not authorized
	}

	// Если удаленная карта была default, назначаем новую default
	if wasDefault {
		setNewDefaultQuery := `
			UPDATE user_payment_methods 
			SET is_default = true, updated_at = NOW()
			WHERE user_id = $1 AND is_active = true
			AND id = (
				SELECT id FROM user_payment_methods 
				WHERE user_id = $1 AND is_active = true 
				ORDER BY created_at ASC 
				LIMIT 1
			)
		`
		_, err = tx.Exec(ctx, setNewDefaultQuery, userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Убираем default у всех карт пользователя
	updateAllQuery := `
		UPDATE user_payment_methods 
		SET is_default = false, updated_at = NOW()
		WHERE user_id = $1 AND is_active = true
	`
	_, err = tx.Exec(ctx, updateAllQuery, userID)
	if err != nil {
		return err
	}

	// Устанавливаем default для выбранной карты
	setDefaultQuery := `
		UPDATE user_payment_methods 
		SET is_default = true, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`
	result, err := tx.Exec(ctx, setDefaultQuery, paymentMethodID, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return err // Not found or not authorized
	}

	return tx.Commit(ctx)
}

func (r *repository) GetDefaultPaymentMethod(ctx context.Context, userID int64) (*models.UserPaymentMethod, error) {
	query := `
		SELECT id, user_id, stripe_payment_method_id, stripe_customer_id,
		       card_last4, card_brand, card_exp_month, card_exp_year,
		       is_default, is_active, created_at, updated_at
		FROM user_payment_methods
		WHERE user_id = $1 AND is_default = true AND is_active = true
		LIMIT 1
	`

	var pm models.UserPaymentMethod
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&pm.ID, &pm.UserID, &pm.StripePaymentMethodID, &pm.StripeCustomerID,
		&pm.CardLast4, &pm.CardBrand, &pm.CardExpMonth, &pm.CardExpYear,
		&pm.IsDefault, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &pm, nil
}

// internal/repository/payment/payments.go
func (r *repository) SavePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (
			user_id, job_id, stripe_payment_intent_id, stripe_payment_method_id,
			stripe_customer_id, amount_cents, currency, status, description
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		payment.UserID,
		payment.JobID,
		payment.StripePaymentIntentID,
		payment.StripePaymentMethodID,
		payment.StripeCustomerID,
		payment.AmountCents,
		payment.Currency,
		payment.Status,
		payment.Description,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	return err
}

func (r *repository) GetPaymentByStripeIntentID(ctx context.Context, stripePaymentIntentID string) (*models.Payment, error) {
	query := `
		SELECT id, user_id, job_id, stripe_payment_intent_id, stripe_payment_method_id,
		       stripe_customer_id, amount_cents, currency, status, description,
		       failure_reason, created_at, updated_at
		FROM payments
		WHERE stripe_payment_intent_id = $1
	`

	var payment models.Payment
	var jobID *int64
	err := r.db.QueryRow(ctx, query, stripePaymentIntentID).Scan(
		&payment.ID, &payment.UserID, &jobID, &payment.StripePaymentIntentID,
		&payment.StripePaymentMethodID, &payment.StripeCustomerID,
		&payment.AmountCents, &payment.Currency, &payment.Status,
		&payment.Description, &payment.FailureReason,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if jobID != nil {
		payment.JobID = *jobID
	}

	return &payment, nil
}

func (r *repository) UpdatePaymentStatus(ctx context.Context, paymentID int64, status, failureReason string) error {
	query := `
		UPDATE payments 
		SET status = $1, failure_reason = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, status, failureReason, paymentID)
	return err
}

func (r *repository) GetUserPayments(ctx context.Context, userID int64, limit, offset int) ([]models.Payment, error) {
	query := `
		SELECT id, user_id, job_id, stripe_payment_intent_id, stripe_payment_method_id,
		       stripe_customer_id, amount_cents, currency, status, description,
		       failure_reason, created_at, updated_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		var jobID *int64
		err := rows.Scan(
			&payment.ID, &payment.UserID, &jobID, &payment.StripePaymentIntentID,
			&payment.StripePaymentMethodID, &payment.StripeCustomerID,
			&payment.AmountCents, &payment.Currency, &payment.Status,
			&payment.Description, &payment.FailureReason,
			&payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if jobID != nil {
			payment.JobID = *jobID
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}
