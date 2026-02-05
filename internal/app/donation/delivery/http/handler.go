package http

import (
	"sakucita/internal/domain"
	"sakucita/internal/infra/midtrans"
	"sakucita/internal/server/middleware"
	"sakucita/pkg/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Handler struct {
	config         config.App
	log            zerolog.Logger
	validator      *validator.Validate
	mw             *middleware.Middleware
	midtransClient *midtrans.MidtransClient
}

func NewHandler(
	config config.App,
	log zerolog.Logger,
	validator *validator.Validate,
	mw *middleware.Middleware,
	midtransClient *midtrans.MidtransClient,
) *Handler {
	return &Handler{
		config,
		log,
		validator,
		mw,
		midtransClient,
	}
}

func (h *Handler) Routes(r fiber.Router) {
	r.Route("/donations", func(router fiber.Router) {
		router.Post("/", h.CreateDonation)
	})
}

func (h *Handler) CreateDonation(c fiber.Ctx) error {
	var req domain.CreateDonationMessageRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(req); err != nil {
		return err
	}

	res, err := h.midtransClient.CreateQRIS(c.RequestCtx(), midtrans.MidtransQRISRequest{
		PaymentType: "qris",
		TransactionDetails: midtrans.MidtransTransactionDetails{
			OrderID:     uuid.New().String(),
			GrossAmount: int64(req.Amount) + int64(750),
		},
		ItemDetails: []midtrans.MidtransItemDetail{
			{
				ID:       uuid.New().String(),
				Price:    int64(req.Amount),
				Quantity: 1,
				Name:     "donation amount",
			},
			{
				ID:       uuid.New().String(),
				Price:    750,
				Quantity: 1,
				Name:     "donation fee",
			},
		},
		CustomerDetails: &midtrans.MidtransCustomerDetails{
			FirstName: req.PayerName,
			Email:     *req.Email,
		},
		QRIS: midtrans.MidtransQRISDetail{
			Acquirer: string(midtrans.QRISGopay),
		},
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Message: "success",
		Data:    res,
	})
}
