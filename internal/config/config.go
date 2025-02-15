package config

import (
	"log"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type dbConfig struct {
	DBConnURL        string `env:"DB_URL"            env-required:"true"`
	DatabasePort     string `env:"DATABASE_PORT"     env-required:"true"`
	DatabaseUser     string `env:"DATABASE_USER"     env-required:"true"`
	DatabasePassword string `env:"DATABASE_PASSWORD" env-required:"true"`
	DatabaseName     string `env:"DATABASE_NAME"     env-required:"true"`
	DatabaseHost     string `env:"DATABASE_HOST"     env-required:"true"`
}

type appConfig struct {
	ServerPort string `env:"SERVER_PORT" env-required:"true"`
}

type Config struct {
	dbConfig
	appConfig
}

func MustInit() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Info("failed to read .env", slog.Any("error", err))
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
