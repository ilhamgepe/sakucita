package security

import (
	"os"
	"testing"

	"sakucita/pkg/config"
	"sakucita/pkg/logger"
)

var testSecurity *Security

func TestMain(m *testing.M) {
	config, err := config.New("../../../config.yaml")
	if err != nil {
		panic(err)
	}

	log := logger.New("testing security", config)
	s := Security{
		config: config,
		log:    log,
	}

	testSecurity = &s
	code := m.Run()

	os.Exit(code)
}
