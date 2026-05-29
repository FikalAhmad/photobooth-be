package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentCanceled  PaymentStatus = "CANCELED"
)

type PaymentMethod string

const (
	PaymentCash         PaymentMethod = "CASH"
	PaymentEWallet      PaymentMethod = "E_WALLET"
	PaymentQRIS         PaymentMethod = "QRIS"
	PaymentBankTransfer PaymentMethod = "BANK_TRANSFER"
)

type Payment struct {
	ID             uuid.UUID
	UserID         string
	SubscriptionID string
	GatewayTrxID   string
	Amount         int
	Status         string
	PaymentType    string
	PaidAt         time.Time
	CreatedAt      time.Time
}
