package security

import (
	"crypto/rsa"

	"sakucita/pkg/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Security struct {
	config    config.App
	log       zerolog.Logger
	rdb       *redis.Client
	activeKID string
	rsaKeys   map[string]*RSAKeys
}

type RSAKeys struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

func NewSecurity(cfg config.App, log zerolog.Logger, rdb *redis.Client) *Security {
	return &Security{
		config: cfg,
		log:    log,
		rdb:    rdb,
	}
}
