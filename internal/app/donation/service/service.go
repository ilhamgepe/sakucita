package service

import (
	"context"

	"sakucita/internal/domain"
	"sakucita/internal/infra/midtrans"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *service) CreateDonation(
	ctx context.Context,
	req CreateDonationCommand,
) (*CreateDonationResult, error) {

	/**
	 * 1. Setup Database Transaction
	 */
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	/**
	 * 2. Validate & Load Required Data
	 */

	creator, err := qtx.GetUserByID(ctx, req.PayeeUserID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return nil, domain.NewAppError(
				fiber.StatusNotFound,
				domain.ErrMsgUserNotFound,
				domain.ErrNotfound,
			)
		}
		s.log.Err(err).Msg("failed to get user by id")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	paymentChannel, err := qtx.GetPaymentChannelByCode(ctx, req.PaymentChannel)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return nil, domain.NewAppError(
				fiber.StatusNotFound,
				domain.ErrMsgPaymentChannelNotFound,
				domain.ErrNotfound,
			)
		}
		s.log.Err(err).Msg("failed to get payment channel")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	fees, err := qtx.GetUserFee(ctx, repository.GetUserFeeParams{
		Userid:           creator.ID,
		Paymentchannelid: paymentChannel.ID,
	})
	if err != nil {
		s.log.Err(err).Msg("failed to get user fee")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 3. Calculate Amount & Fees
	 */

	donationAmount := int64(req.Amount)

	gatewayFee := (donationAmount*fees.GatewayFeePercentage)/10000 +
		fees.GatewayFeeFixed

	platformFee := (donationAmount*fees.PlatformFeePercentage)/10000 +
		fees.PlatformFeeFixed

	totalFeeFixed := fees.GatewayFeeFixed + fees.PlatformFeeFixed
	totalFeePercentage := fees.GatewayFeePercentage + fees.PlatformFeePercentage
	totalFeeAmount := gatewayFee + platformFee

	grossAmount := donationAmount + gatewayFee
	netAmount := grossAmount - totalFeeAmount

	/**
	 * 4. Create Donation Message
	 */

	donationMsgID, err := utils.GenerateUUIDV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate donation message id")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	donationMsg, err := qtx.CreateDonationMessage(ctx,
		repository.CreateDonationMessageParams{
			ID:                donationMsgID,
			PayeeUserID:       req.PayeeUserID,
			PayerUserID:       utils.StringToPgTypeUUID(req.PayerUserID.String()),
			PayerName:         req.PayerName,
			Email:             req.Email,
			Message:           req.Message,
			MediaType:         repository.DonationMediaType(req.MediaType),
			MediaUrl:          utils.StringPtrToPgTypeText(req.MediaURL),
			MediaStartSeconds: utils.Int32PtrToPgTypeInt4(req.MediaStartSeconds),
			MaxPlaySeconds:    utils.MaxPlayedSeconds(int32(req.Amount), 500),
			PricePerSecond:    pgtype.Int8{Int64: 500, Valid: true},
			Amount:            donationAmount,
			Currency:          "IDR",
			Meta:              domain.JSONB{},
		},
	)
	if err != nil {
		s.log.Err(err).Msg("failed to create donation message")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 5. Call Payment Gateway
	 */

	qrisResult, err := s.midtransClient.CreateQRIS(
		ctx,
		donationAmount,
		req.PayerName,
		req.Email,
		gatewayFee,
	)
	if err != nil {
		s.log.Err(err).Msg("failed to create qris midtrans")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			"failed to create qris",
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 6. Create Transaction Record
	 */

	transactionID, err := uuid.NewV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate transaction id")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	transaction, err := qtx.CreateTransaction(ctx,
		repository.CreateTransactionParams{
			ID:                    transactionID,
			DonationMessageID:     donationMsg.ID,
			PaymentChannelID:      paymentChannel.ID,
			PayeeUserID:           req.PayeeUserID,
			PayerUserID:           utils.StringToPgTypeUUID(req.PayerUserID.String()),
			Amount:                grossAmount,
			GatewayFeeFixed:       fees.GatewayFeeFixed,
			GatewayFeePercentage:  fees.GatewayFeePercentage,
			GatewayFeeAmount:      gatewayFee,
			PlatformFeeFixed:      fees.PlatformFeeFixed,
			PlatformFeePercentage: fees.PlatformFeePercentage,
			PlatformFeeAmount:     platformFee,
			FeeFixed:              totalFeeFixed,
			FeePercentage:         totalFeePercentage,
			FeeAmount:             totalFeeAmount,
			NetAmount:             netAmount,
			Currency:              "IDR",
			Status:                "PENDING",
			ExternalReference:     utils.StringToPgTypeText(qrisResult.TransactionID),
		},
	)
	if err != nil {
		s.log.Err(err).Msg("failed to create transaction")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 7. Commit Transaction
	 */

	if err := tx.Commit(ctx); err != nil {
		s.log.Err(err).Msg("failed to commit transaction")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 8. Build Response
	 */

	return &CreateDonationResult{
		TransactionID: transactionID.String(),
		Amount:        grossAmount,
		Currency:      transaction.Currency,
		Status:        string(transaction.Status),
		QrString:      qrisResult.QRString,
	}, nil
}

// func (s *service) CreateDonation(ctx context.Context, req CreateDonationCommand) (*CreateDonationResult, error) {
// 	// TODO
// 	// setup db transaction
// 	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}
// 	defer func() {
// 		_ = tx.Rollback(ctx)
// 	}()

// 	qtx := s.q.WithTx(tx)

// 	// get payee user or creator
// 	creator, err := qtx.GetUserByID(ctx, req.PayeeUserID)
// 	if err != nil {
// 		if utils.IsNotFoundError(err) {
// 			return nil, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgUserNotFound, domain.ErrNotfound)
// 		}

// 		s.log.Err(err).Msg("failed to get user by id")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	// get payment channel
// 	paymentChannel, err := qtx.GetPaymentChannelByCode(ctx, req.PaymentChannel)
// 	if err != nil {
// 		if utils.IsNotFoundError(err) {
// 			return nil, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgPaymentChannelNotFound, domain.ErrNotfound)
// 		}

// 		s.log.Err(err).Msg("failed to get payment channel by code")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	// get fee config
// 	fees, err := qtx.GetUserFee(ctx, repository.GetUserFeeParams{
// 		Userid:           creator.ID,
// 		Paymentchannelid: paymentChannel.ID,
// 	})
// 	if err != nil {
// 		s.log.Err(err).Msg("failed to get user fee")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	// create donation message
// 	donationMsgID, err := utils.GenerateUUIDV7()
// 	if err != nil {
// 		s.log.Err(err).Msg("failed to generate uuid v7 for donation message")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	donationMsgResult, err := qtx.CreateDonationMessage(ctx, repository.CreateDonationMessageParams{
// 		ID:                donationMsgID,
// 		PayeeUserID:       req.PayeeUserID,
// 		PayerUserID:       utils.StringToPgTypeUUID(req.PayerUserID.String()),
// 		PayerName:         req.PayerName,
// 		Email:             req.Email,
// 		Message:           req.Message,
// 		MediaType:         repository.DonationMediaType(req.MediaType),
// 		MediaUrl:          utils.StringPtrToPgTypeText(req.MediaURL),
// 		MediaStartSeconds: utils.Int32PtrToPgTypeInt4(req.MediaStartSeconds),
// 		MaxPlaySeconds:    utils.MaxPlayedSeconds(int32(req.Amount), 500),
// 		PricePerSecond:    pgtype.Int8{Int64: 500, Valid: true},
// 		Amount:            int64(req.Amount),
// 		Currency:          "IDR",
// 		Meta:              domain.JSONB{},
// 	})
// 	if err != nil {
// 		s.log.Err(err).Msg("failed to create donation message")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	// gateway fee
// 	gatewayFee := (int64(req.Amount) * fees.GatewayFeePercentage / 10000) + fees.GatewayFeeFixed

// 	// create qris midtrans
// 	qrisResult, err := s.midtransClient.CreateQRIS(ctx, int64(req.Amount), req.PayerName, req.Email, gatewayFee)
// 	if err != nil {
// 		s.log.Err(err).Msg("failed to create qris midtrans")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, "failed to create qris", domain.ErrInternalServerError)
// 	}
// 	grossAmount, err := utils.ParseRupiahAmount(qrisResult.GrossAmount)
// 	if err != nil {
// 		s.log.Err(err).Msg("failed to parse rupiah amount")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	// platform fee
// 	platformFee := (req.Amount * fees.PlatformFeePercentage / 10000) + fees.PlatformFeeFixed

// 	// total fee
// 	totalFeeFixed := fees.GatewayFeeFixed + fees.PlatformFeeFixed
// 	totalFeePercentage := fees.GatewayFeePercentage + fees.PlatformFeePercentage
// 	totalFeeAmount := gatewayFee + platformFee

// 	// create transaction and get transaction id from midtrans for external reference
// 	transactionID, err := uuid.NewV7()
// 	if err != nil {
// 		s.log.Err(err).Msg("failed to generate uuid v7 for transaction id")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	transactionResult, err := qtx.CreateTransaction(ctx, repository.CreateTransactionParams{
// 		ID:                    transactionID,
// 		DonationMessageID:     donationMsgResult.ID,
// 		PaymentChannelID:      paymentChannel.ID,
// 		PayeeUserID:           req.PayeeUserID,
// 		PayerUserID:           utils.StringToPgTypeUUID(req.PayerUserID.String()),
// 		Amount:                grossAmount,
// 		GatewayFeeFixed:       fees.GatewayFeeFixed,
// 		GatewayFeePercentage:  fees.GatewayFeePercentage,
// 		GatewayFeeAmount:      gatewayFee,
// 		PlatformFeeFixed:      fees.PlatformFeeFixed,
// 		PlatformFeePercentage: fees.PlatformFeePercentage,
// 		PlatformFeeAmount:     platformFee,
// 		FeeFixed:              totalFeeFixed,
// 		FeePercentage:         totalFeePercentage,
// 		FeeAmount:             totalFeeAmount,
// 		NetAmount:             grossAmount - totalFeeAmount,
// 		Currency:              "IDR",
// 		Status:                "PENDING",
// 		ExternalReference:     utils.StringToPgTypeText(qrisResult.TransactionID),
// 	})
// 	if err != nil {
// 		s.log.Err(err).Msg("failed to create transaction")
// 		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
// 	}

// 	// kasih response qr string aja dulu biar fe yang buat image qr nya.

// 	return &CreateDonationResult{
// 		TransactionID: transactionID.String(),
// 		Amount:        grossAmount,
// 		Currency:      transactionResult.Currency,
// 		Status:        transactionResult.Currency,
// 		QrString:      qrisResult.QRString,
// 	}, nil
// }
