package main

import (
	"context"

	authService "sakucita/internal/app/auth/service"
	"sakucita/internal/database"
	"sakucita/internal/database/repository"
	"sakucita/internal/domain"
	"sakucita/internal/server"
	"sakucita/pkg/config"
	"sakucita/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func main() {
	cfg := configProvider()
	log := loggerProvider(cfg)
	databases := databaseProvider(cfg, log)

	queries := repository.New(databases.postgres)

	services := serviceProvider(log, databases, queries)

	serverHttp := ServerHTTPProvider(cfg, log, services)

	serverHttp.Start()
}

// service provider
type services struct {
	authService domain.AuthService
}

func serviceProvider(log zerolog.Logger, databases *databases, queries *repository.Queries) *services {
	return &services{
		authService: authService.NewService(databases.postgres, queries, log),
	}
}

// database provider
type databases struct {
	postgres *pgxpool.Pool
	redis    *redis.Client
}

func databaseProvider(cfg config.App, log zerolog.Logger) *databases {
	pg, err := database.NewDB(context.Background(), cfg, log)
	if err != nil {
		panic(err)
	}

	redis, err := database.NewRedisClient(cfg, log)
	if err != nil {
		panic(err)
	}
	return &databases{
		postgres: pg,
		redis:    redis,
	}
}

// config provider
func configProvider() config.App {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	return cfg
}

// logger provider
func loggerProvider(cfg config.App) zerolog.Logger {
	return logger.New("sakucita", cfg)
}

// server provider
func ServerHTTPProvider(cfg config.App, log zerolog.Logger, services *services) *server.Server {
	return server.NewServer(
		cfg, log, services.authService,
	)
}
