package http

import (
	"sakucita/internal/app/auth/service"
	"sakucita/internal/domain"
	"sakucita/internal/dto"
	"sakucita/internal/server/middleware"
	"sakucita/internal/server/security"
	"sakucita/internal/shared/utils"
	"sakucita/pkg/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type Handler struct {
	config      config.App
	log         zerolog.Logger
	validator   *validator.Validate
	authService service.AuthService
	mw          *middleware.Middleware
}

func NewHandler(config config.App, log zerolog.Logger, validator *validator.Validate, authService service.AuthService, mw *middleware.Middleware) *Handler {
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
	r.Post("/auth/login", h.mw.LoginLimiter, h.loginLocal)

	r.Route("/auth", func(router fiber.Router) {
		router.Use(h.mw.WithAuth)
		router.Get("/me", h.me)
		router.Get("/auth/refresh", h.refreshToken)
	})
}

func (h *Handler) refreshToken(c fiber.Ctx) error {
	claims, ok := c.Locals(domain.CtxUserIDKey).(security.TokenClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: domain.ErrUnauthorized.Error(),
			Errors:  "invalid token claims",
		})
	}

	res, err := h.authService.RefreshToken(c.RequestCtx(), service.RefreshCommand{
		Claims:     claims,
		ClientInfo: utils.ExtractClientInfo(c),
	})
	if err != nil {
		return err
	}

	return c.JSON(dto.Response{
		Message: "success",
		Data:    res,
	})
}

func (h *Handler) me(c fiber.Ctx) error {
	claims, ok := c.Locals(domain.CtxUserIDKey).(security.TokenClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: domain.ErrUnauthorized.Error(),
			Errors:  "invalid token claims",
		})
	}

	user, err := h.authService.Me(c.RequestCtx(), claims.UserID)
	if err != nil {
		return err
	}

	return c.JSON(dto.Response{
		Message: "success",
		Data:    user,
	})
}

func (h *Handler) loginLocal(c fiber.Ctx) error {
	var req service.LoginLocalCommand
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	req.ClientInfo = utils.ExtractClientInfo(c)

	res, err := h.authService.LoginLocal(c.RequestCtx(), req)
	if err != nil {
		return err
	}
	return c.JSON(dto.Response{
		Message: "success",
		Data:    res,
	})
}

func (h *Handler) registerLocal(c fiber.Ctx) error {
	var req service.RegisterCommand
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(req); err != nil {
		return err
	}

	if err := h.authService.RegisterLocal(c.RequestCtx(), req); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Message: "register success",
	})
}

// fiber:context-methods migrated
