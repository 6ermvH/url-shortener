package main

import (
	"fmt"
	"log"

	"github.com/6ermvH/url-shortener/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log.Printf("Url-Shortener service: addr=%s storage=%s", cfg.Addr, cfg.Storage)

	return nil
}
