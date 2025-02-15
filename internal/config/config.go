package config

import (
	"log"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConnURL string `env:"DB_URL" env-required:"true"`
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
