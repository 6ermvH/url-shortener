package config

import (
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Addr    string `validate:"required,hostname_port"`
	Storage string `validate:"required,oneof=postgres memory"`
	DSN     string `validate:"required_if=Storage postgres"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Addr, "addr", ":8080", "server address")
	flag.StringVar(&cfg.Storage, "storage", "memory", "storage type: postgres|memory")
	flag.StringVar(&cfg.DSN, "dsn", "", "postgres DSN (required when storage=postgres)")
	flag.Parse()

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}
