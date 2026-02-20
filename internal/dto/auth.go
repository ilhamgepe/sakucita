package dto

import (
	"sakucita/internal/domain"
	"sakucita/internal/server/security"
)

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
	ClientInfo security.ClientInfo
}

type LoginResponse struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

type RefreshRequest struct {
	Claims     security.TokenClaims
	ClientInfo security.ClientInfo
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
