package domain

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterLocal(ctx context.Context, req RegisterRequest) error
	LoginLocal(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Me(ctx context.Context, userID uuid.UUID) (*UserWithRoles, error)
}

type RegisterRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Phone    string `json:"phone" form:"phone" validate:"required,min=10,max=15"`
	Name     string `json:"name" form:"name" validate:"required,min=2,max=30"`
	Nickname string `json:"nickname" form:"nickname" validate:"required,min=2,max=30"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=200"`
}

type LoginRequest struct {
	Email      string `json:"email" form:"email" validate:"required,email"`
	Password   string `json:"password" form:"password" validate:"required,min=8,max=200"`
	ClientInfo ClientInfo
}

type ClientInfo struct {
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceName string `json:"device_name"`
}

type LoginResponse struct {
	User         UserWithRoles `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}
