package middleware

import (
	"sakucita/internal/server/security"

	"github.com/rs/zerolog"
)

type Middleware struct {
	log      zerolog.Logger
	security *security.Security
}

func NewMiddleware(log zerolog.Logger, security *security.Security) *Middleware {
	return &Middleware{log, security}
}
