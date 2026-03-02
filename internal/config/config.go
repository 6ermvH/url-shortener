package config

import "flag"

type Config struct {
	Addr    string
	Storage string
	DSN     string
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Addr, "addr", ":8080", "server address")
	flag.StringVar(&cfg.Storage, "storage", "memory", "storage type: postgres|memory")
	flag.StringVar(&cfg.DSN, "dsn", "", "postgres DSN (required when storage=postgres)")
	flag.Parse()

	return cfg
}
