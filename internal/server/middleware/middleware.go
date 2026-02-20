package middleware

import (
	"sakucita/internal/app/auth/service"
	"sakucita/internal/server/security"

	"github.com/rs/zerolog"
)

type Middleware struct {
	log         zerolog.Logger
	security    *security.Security
	authService service.AuthService
}

func NewMiddleware(log zerolog.Logger, security *security.Security, authService service.AuthService) *Middleware {
	return &Middleware{log, security, authService}
}
