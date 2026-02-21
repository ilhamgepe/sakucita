package midtrans

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (c *midtransClient) CreateQRIS(
	ctx context.Context,
	amount int64,
	payerName, payerEmail string,
	midtransQRISFee int64,

) (*MidtransQRISResponse, error) {
	req := MidtransQRISRequest{
		PaymentType: "qris",
		TransactionDetails: MidtransTransactionDetails{
			OrderID:     uuid.New().String(),
			GrossAmount: int64(amount) + int64(midtransQRISFee),
		},
		ItemDetails: []MidtransItemDetail{
			{
				ID:       uuid.New().String(),
				Price:    int64(amount),
				Quantity: 1,
				Name:     "donation amount",
			},
			{
				ID:       uuid.New().String(),
				Price:    midtransQRISFee,
				Quantity: 1,
				Name:     "payment gatewway fee",
			},
		},
		CustomerDetails: &MidtransCustomerDetails{
			FirstName: payerName,
			Email:     payerEmail,
		},
		QRIS: MidtransQRISDetail{
			Acquirer: string(QRISGopay),
		},
	}
	var resp MidtransQRISResponse
	r, err := c.http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&resp).
		Post("/v2/charge")
	if err != nil {
		c.log.Err(err).Msg("error http request")
		return nil, err
	}

	if r.IsError() {
		c.log.Error().Msgf("midtrans error: %s", r.String())
		return nil, fmt.Errorf("midtrans error: %s", r.String())
	}

	return &resp, nil
}
