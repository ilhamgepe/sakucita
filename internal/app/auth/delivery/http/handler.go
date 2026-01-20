package http

import (
	"errors"

	"sakucita/internal/domain"
	"sakucita/internal/server/middleware"
	"sakucita/internal/shared/utils"
	"sakucita/pkg/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	config      config.App
	log         zerolog.Logger
	validator   *validator.Validate
	authService domain.AuthService
	mw          *middleware.Middleware
}

func NewHandler(config config.App, log zerolog.Logger, validator *validator.Validate, authService domain.AuthService, mw *middleware.Middleware) *Handler {
	return &Handler{
		config,
		log,
		validator,
		authService,
		mw,
	}
}

func (h *Handler) Routes(r fiber.Router) {
	r.Post("/auth/register", h.registerLocal)
	r.Post("/auth/login", h.loginLocal)

	r.Get("/auth/me", h.mw.WithAuth, h.me)
}

func (h *Handler) me(c *fiber.Ctx) error {
	claims, ok := c.Locals(domain.CtxUserIDKey).(domain.TokenClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
			Message: domain.ErrUnauthorized.Error(),
			Errors:  "invalid token claims",
		})
	}

	user, err := h.authService.Me(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.JSON(domain.Response{
		Message: "success",
		Data:    user,
	})
}

func (h *Handler) loginLocal(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	req.ClientInfo = utils.ExtractClientInfo(c)

	res, err := h.authService.LoginLocal(c.Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(domain.Response{
		Message: "success",
		Data:    res,
	})
}

func (h *Handler) registerLocal(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(req); err != nil {
		return err
	}

	if err := h.authService.RegisterLocal(c.Context(), req); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) ||
			errors.Is(err, domain.ErrPhoneAlreadyExists) ||
			errors.Is(err, domain.ErrNicknameAlreadyExists) {

			return c.Status(fiber.StatusConflict).JSON(domain.ErrorResponse{
				Message: "register failed",
				Errors:  err.Error(),
			})
		}
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Message: "register success",
	})
}
