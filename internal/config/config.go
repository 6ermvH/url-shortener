package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Addr    string `validate:"required,hostname_port"`
	Storage string `validate:"required,oneof=postgres memory"`
	DSN     string `validate:"required_if=Storage postgres"`
}

func Load() (*Config, error) {
	cfg := &Config{ //nolint:exhaustruct
		Addr:    envOrDefault("ADDR", ":8080"),
		Storage: envOrDefault("STORAGE", "memory"),
		DSN:     envOrDefault("DSN", ""),
	}

	flag.StringVar(&cfg.Addr, "addr", cfg.Addr, "server address")
	flag.StringVar(&cfg.Storage, "storage", cfg.Storage, "storage type: postgres|memory")
	flag.StringVar(&cfg.DSN, "dsn", cfg.DSN, "postgres DSN (required when storage=postgres)")
	flag.Parse()

	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}
