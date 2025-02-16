package config

import (
	"log"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabasePort     string        `env:"DATABASE_PORT"     env-required:"true"`
	DatabaseUser     string        `env:"DATABASE_USER"     env-required:"true"`
	DatabasePassword string        `env:"DATABASE_PASSWORD" env-required:"true"`
	DatabaseName     string        `env:"DATABASE_NAME"     env-required:"true"`
	DatabaseHost     string        `env:"DATABASE_HOST"     env-required:"true"`
	ServerPort       string        `env:"SERVER_PORT"       env-required:"true"`
	JWTSecret        string        `env:"JWT_SECRET"        env-required:"true"`
	AccessTokenTTL   time.Duration `env:"ACCESS_TOKEN_TTL"  env-required:"true"`
	PublicRoutes     []string
}

func MustInit(envPath string) Config {
	var cfg Config

	cfg.PublicRoutes = newPublicRoutes()

	if err := cleanenv.ReadConfig(envPath, &cfg); err != nil {
		slog.Info("failed to read .env", slog.Any("error", err))
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func newPublicRoutes() []string {
	return []string{
		"/api/auth",
	}
}
