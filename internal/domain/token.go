package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   []Role    `json:"role"`
	jwt.RegisteredClaims
}
