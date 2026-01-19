package logger

import (
	"os"
	"time"

	"sakucita/pkg/config"

	"github.com/rs/zerolog"
)

func New(service string, config config.App) zerolog.Logger {
	if config.Server.Env == "development" {
		return consoleLogger(service)
	} else {
		return jsonLogger(service)
	}
}

func consoleLogger(service string) zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatTimestamp: func(i interface{}) string {
			if t, ok := i.(string); ok {
				return t
			}
			return time.Now().Format(time.RFC3339)
		},
	}

	l := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("service", service).
		Logger()

	return l
}

func jsonLogger(service string) zerolog.Logger {
	l := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("service", service).
		Logger()

	return l
}
