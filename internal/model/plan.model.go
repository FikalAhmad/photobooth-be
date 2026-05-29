package model

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID             uuid.UUID //0129312093
	Name           string    // Starter Strip, Arcade Pass Pro, Studio Elite VIP
	PlanType       string    // FREE, PRO, VIP
	Price          int       // 0, 29000, 99000
	Billing        string    // no, month, month
	IncludeAI      bool      // true, true, true
	IncludeStorage bool      // false, true, true
	CreatedAt      time.Time
}
