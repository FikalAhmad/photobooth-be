package repository

import (
	"database/sql"
	"photobooth-be/internal/model"
)

type PaymentRepository interface {
	CreatePayment(payment *model.Payment) error
	GetPayment() (*model.Payment, error)
	UpdatePayment(id string, payment *model.Payment) error
	DeletePayment(id string) error
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) CreatePayment(payment *model.Payment) error {
	query := `
    INSERT INTO payments (user_id, subscription_id, gateway_trx_id, amount, status, payment_type, paid_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id
  `

	err := r.db.QueryRow(
		query,
		payment.UserID,
		payment.SubscriptionID,
		payment.GatewayTrxID,
		payment.Amount,
		payment.Status,
		payment.PaymentType,
		payment.PaidAt).Scan(&payment.UserID)

	return err
}

func (r *paymentRepository) GetPayment() (*model.Payment, error) {
	query := `
		SELECT id, user_id, subscription_id, gateway_trx_id, amount, status, payment_type, paid_at
		FROM payments
	`

	var payment model.Payment

	err := r.db.QueryRow(query).Scan(
		&payment.UserID,
		&payment.SubscriptionID,
		&payment.GatewayTrxID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentType,
		&payment.PaidAt)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) UpdatePayment(id string, payment *model.Payment) error {
	query := `
    INSERT INTO payments (user_id, subscription_id, gateway_trx_id, amount, status, payment_type, paid_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id
  `

	err := r.db.QueryRow(
		query,
		payment.UserID,
		payment.SubscriptionID,
		payment.GatewayTrxID,
		payment.Amount,
		payment.Status,
		payment.PaymentType,
		payment.PaidAt).Scan(&payment.UserID)

	return err
}

func (r *paymentRepository) DeletePayment(id string) error {
	query := `
    INSERT INTO payments (user_id, subscription_id, gateway_trx_id, amount, status, payment_type, paid_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id
  `

	_, err := r.db.Exec(query)

	return err
}
