package http

import (
	"sakucita/internal/app/donation/service"
	"sakucita/internal/domain"
	"sakucita/internal/dto"
	"sakucita/internal/server/middleware"
	"sakucita/pkg/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type Handler struct {
	config    config.App
	log       zerolog.Logger
	validator *validator.Validate
	mw        *middleware.Middleware
	service   service.DonationService
}

func NewHandler(
	config config.App,
	log zerolog.Logger,
	validator *validator.Validate,
	mw *middleware.Middleware,
	service service.DonationService,
) *Handler {
	return &Handler{
		config,
		log,
		validator,
		mw,
		service,
	}
}

func (h *Handler) Routes(r fiber.Router) {
	r.Route("/donations", func(router fiber.Router) {
		router.Post("/", h.CreateDonation)
	})
}

func (h *Handler) CreateDonation(c fiber.Ctx) error {
	var req dto.CreateDonationRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return domain.NewAppError(fiber.StatusBadRequest, "failed", err)
	}
	res, err := h.service.CreateDonation(c.RequestCtx(), service.CreateDonationCommand{
		PayeeUserID:       req.PayeeUserID,
		PayerUserID:       req.PayerUserID,
		PayerName:         req.PayerName,
		Email:             req.Email,
		Message:           req.Message,
		MediaType:         req.MediaType,
		MediaURL:          req.MediaURL,
		MediaStartSeconds: req.MediaStartSeconds,
		Amount:            req.Amount,
		PaymentChannel:    req.PaymentChannel,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Message: "success",
		Data:    res,
	})
}
