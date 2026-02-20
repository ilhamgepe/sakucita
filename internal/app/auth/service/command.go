package service

import (
	"sakucita/internal/domain"
	"sakucita/internal/server/security"
)

type RegisterCommand struct {
	Email    string
	Phone    string
	Name     string
	Nickname string
	Password string
}

type LoginLocalCommand struct {
	Email    string
	Password string

	ClientInfo security.ClientInfo
}

type LoginResult struct {
	User         domain.UserWithRoles
	AccessToken  string
	RefreshToken string
}

type RefreshCommand struct {
	Claims     security.TokenClaims
	ClientInfo security.ClientInfo
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}
