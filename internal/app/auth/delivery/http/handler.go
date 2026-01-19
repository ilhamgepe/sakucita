package http

import (
	"errors"

	"sakucita/internal/domain"
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
}

func NewHandler(config config.App, log zerolog.Logger, validator *validator.Validate, authService domain.AuthService) *Handler {
	return &Handler{
		config,
		log,
		validator,
		authService,
	}
}

func (h *Handler) Routes(r fiber.Router) {
	r.Post("/auth/register", h.registerLocal)
	r.Post("/auth/login", h.loginLocal)
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
