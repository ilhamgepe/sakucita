package security

import (
	"sakucita/pkg/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Security struct {
	config config.App
	log    zerolog.Logger
	rdb    *redis.Client
}

func NewSecurity(cfg config.App, log zerolog.Logger, rdb *redis.Client) *Security {
	return &Security{
		config: cfg,
		log:    log,
		rdb:    rdb,
	}
}
