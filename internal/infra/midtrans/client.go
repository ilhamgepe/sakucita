package midtrans

import (
	"context"
	"time"

	"sakucita/pkg/config"

	"github.com/rs/zerolog"
	"resty.dev/v3"
)

type midtransClient struct {
	http *resty.Client
	log  zerolog.Logger
}

type MidtransClient interface {
	CreateQRIS(ctx context.Context, amount int64, payerName, payerEmail string, midtransQRISFee int64) (*MidtransQRISResponse, error)
}

func NewMidtransClient(config config.App, log zerolog.Logger) MidtransClient {
	c := resty.New().
		SetBaseURL(config.Midtrans.BaseURL).
		SetBasicAuth(config.Midtrans.ServerKey, "").
		SetTimeout(15*time.Second).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	return &midtransClient{http: c, log: log}
}
