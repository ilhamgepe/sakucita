package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	PROVIDERLOCAL  string = "local"
	PROVIDERGOOGLE string = "google"
	PROVIDERAPPLE  string = "apple"
)

type AuthIdentity struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	Provider     string     `json:"provider"`
	ProviderID   string     `json:"provider_id"`
	PasswordHash *string    `json:"password_hash,omitempty"`
	TotpSecret   *string    `json:"totp_secret,omitempty"`
	TotpEnabled  bool       `json:"totp_enabled"`
	Meta         JSONB      `json:"meta"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
