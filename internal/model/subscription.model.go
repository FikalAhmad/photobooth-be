package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID
	UserID    string
	PlanID    string
	Status    string    //
	StartedAt time.Time // waktu saat bayar
	ExpiredAt time.Time // sama tanggalya  + 1 bulan
	CreatedAt time.Time
}
