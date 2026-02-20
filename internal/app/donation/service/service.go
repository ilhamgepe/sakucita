package service

import (
	"context"

	"sakucita/internal/infra/midtrans"
	"sakucita/internal/infra/postgres/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type service struct {
	db             *pgxpool.Pool
	q              *repository.Queries
	log            zerolog.Logger
	midtransClient midtrans.MidtransClient
}

type DonationService interface {
	CreateDonation(ctx context.Context, req CreateDonationCommand) (*CreateDonationResult, error)
}

func NewService(
	db *pgxpool.Pool,
	q *repository.Queries,
	log zerolog.Logger,
	midtransClient midtrans.MidtransClient,
) DonationService {
	return &service{db, q, log, midtransClient}
}

func (s *service) CreateDonation(ctx context.Context, req CreateDonationCommand) (*CreateDonationResult, error) {
	// // TODO
	// // setup db transaction
	// tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	// if err != nil {
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }
	// defer func() {
	// 	_ = tx.Rollback(ctx)
	// }()

	// qtx := s.q.WithTx(tx)

	// // get payee user or creator
	// creator, err := qtx.GetUserByID(ctx, req.PayeeUserID)
	// if err != nil {
	// 	if utils.IsNotFoundError(err) {
	// 		return nil, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgUserNotFound, domain.ErrNotfound)
	// 	}

	// 	s.log.Err(err).Msg("failed to get user by id")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }

	// // get payment channel
	// paymentChannel, err := s.q.GetPaymentChannelByCode(ctx, req.PaymentChannel)
	// if err != nil {
	// 	if utils.IsNotFoundError(err) {
	// 		return nil, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgPaymentChannelNotFound, domain.ErrNotfound)
	// 	}

	// 	s.log.Err(err).Msg("failed to get payment channel by code")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }

	// // get fee config
	// // ! ini rada bingung nih gimana enak nya. karna harus hindari manggil datbaase 2x dalam 1 transaksi
	// // mungkin buat fee service nih enaknya, jadi jangan sampe donation tau cara hitung fee karna domain nya fee service
	// var platformFee struct{ Percentage, Fixed int64 }
	// if creator.CustomFee {
	// 	customFee, err := s.q.GetUserFeeOverrideByUserID(ctx, repository.GetUserFeeOverrideByUserIDParams{UserID: creator.ID, PaymentChannelID: paymentChannel.ID})
	// 	if err != nil {
	// 		if utils.IsNotFoundError(err) {
	// 			return nil, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgUserFeeOverrideNotFound, domain.ErrNotfound)
	// 		}
	// 	}
	// }

	// // create donation message
	// donationMsgID, err := uuid.NewV7()
	// if err != nil {
	// 	s.log.Err(err).Msg("failed to generate uuid v7 for donation message")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }

	// pricePerSecond := 500
	// maxPlaySeconds := int32(req.Amount / int32(pricePerSecond))

	// donationMsgResult, err := s.q.CreateDonationMessage(ctx, repository.CreateDonationMessageParams{
	// 	ID:                donationMsgID,
	// 	PayeeUserID:       req.PayeeUserID,
	// 	PayerUserID:       utils.StringToPgTypeUUID(req.PayerUserID.String()),
	// 	PayerName:         req.PayerName,
	// 	Email:             req.Email,
	// 	Message:           req.Message,
	// 	MediaType:         repository.DonationMediaType(req.MediaType),
	// 	MediaUrl:          utils.StringPtrToPgTypeText(req.MediaURL),
	// 	MediaStartSeconds: utils.Int32PtrToPgTypeInt4(req.MediaStartSeconds),
	// 	MaxPlaySeconds:    utils.Int32ToPgTypeInt4(maxPlaySeconds),
	// 	PricePerSecond:    pgtype.Int8{Int64: int64(pricePerSecond), Valid: true},
	// 	Amount:            int64(req.Amount),
	// 	Currency:          "IDR",
	// 	Meta:              domain.JSONB{},
	// })
	// if err != nil {
	// 	s.log.Err(err).Msg("failed to create donation message")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }

	// // create qris midtrans
	// qrisResult, err := s.midtransClient.CreateQRIS(ctx, int64(req.Amount), req.PayerName, req.Email)
	// if err != nil {
	// 	s.log.Err(err).Msg("failed to create qris midtrans")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, "failed to create qris", domain.ErrInternalServerError)
	// }
	// grossAmount, err := utils.ParseRupiahAmount(qrisResult.GrossAmount)
	// if err != nil {
	// 	s.log.Err(err).Msg("failed to parse rupiah amount")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }

	// // create transaction and get transaction id from midtrans for external reference
	// transactionID, err := uuid.NewV7()
	// if err != nil {
	// 	s.log.Err(err).Msg("failed to generate uuid v7 for transaction id")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }

	// var fee pgtype.Numeric
	// if err := fee.Scan("0.007"); err != nil {
	// 	s.log.Err(err).Msg("failed to scan fee")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }
	// transactionResult, err := s.q.CreateTransaction(ctx, repository.CreateTransactionParams{
	// 	ID:                transactionID,
	// 	DonationMessageID: donationMsgResult.ID,
	// 	PayeeUserID:       req.PayeeUserID,
	// 	PayerUserID:       utils.StringToPgTypeUUID(req.PayerUserID.String()),
	// 	Amount:            grossAmount,
	// 	FeeFixed:          750,
	// 	FeePercentage:     fee,
	// 	FeeAmount:         750 + int64(fee),
	// })
	// if err != nil {
	// 	s.log.Err(err).Msg("failed to create transaction")
	// 	return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	// }

	// // kasih response qr string aja dulu biar fe yang buat image qr nya.

	return &CreateDonationResult{}, nil
}
