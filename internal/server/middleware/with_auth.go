package middleware

import (
	"strings"

	"sakucita/internal/domain"
	"sakucita/internal/dto"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) WithAuth(c fiber.Ctx) error {
	bearer := c.Get("Authorization")
	if bearer == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: domain.ErrUnauthorized.Error(),
			Errors:  "bearer token required",
		})
	}

	if !strings.HasPrefix(bearer, "Bearer ") {
		c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: domain.ErrUnauthorized.Error(),
			Errors:  "invalid token format",
		})
	}

	token := bearer[len("Bearer "):]
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: domain.ErrUnauthorized.Error(),
			Errors:  "token required",
		})
	}

	claims, err := m.security.VerifyToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: domain.ErrUnauthorized.Error(),
			Errors:  err.Error(),
		})
	}

	c.Locals(domain.CtxUserIDKey, claims)

	return c.Next()
}
