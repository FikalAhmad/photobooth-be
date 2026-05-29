package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	GoogleID     string
	Email        string
	Name         string
	AvatarURL    string
	RefreshToken string
	AICredits    int
	CreatedAt    time.Time
}
