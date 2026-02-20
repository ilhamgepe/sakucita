package security

import (
	"sakucita/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenClaims struct {
	UserID uuid.UUID     `json:"user_id"`
	Role   []domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type RoleClaim struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type ClientInfo struct {
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceName string `json:"device_name"`
}
