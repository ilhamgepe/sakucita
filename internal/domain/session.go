package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	DeviceID       string    `json:"device_id"`
	RefreshTokenID uuid.UUID `json:"refresh_token_id"`
	ExpiresAt      time.Time `json:"expires_at"`
	Revoked        bool      `json:"revoked"`
	Meta           JSONB     `json:"meta,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	LastUsedAt     time.Time `json:"last_used_at"`
}
