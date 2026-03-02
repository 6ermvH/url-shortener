package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/6ermvH/url-shortener/cmd/internal/migrator"
	"github.com/6ermvH/url-shortener/internal/config"
	"github.com/6ermvH/url-shortener/internal/handler"
	"github.com/6ermvH/url-shortener/internal/repository"
	"github.com/6ermvH/url-shortener/internal/repository/memory"
	"github.com/6ermvH/url-shortener/internal/repository/postgres"
	"github.com/6ermvH/url-shortener/internal/service"
)

var errUnknownStorage = errors.New("unknown storage type")

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	repo, err := newRepository(cfg)
	if err != nil {
		return fmt.Errorf("init repository: %w", err)
	}

	svc := service.New(repo)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.Shorten)
	mux.HandleFunc("GET /{short}", h.Resolve)

	const readHeaderTimeout = 5 * time.Second

	//nolint:exhaustruct
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Printf("starting server: addr=%s storage=%s", cfg.Addr, cfg.Storage)

	serveErr := server.ListenAndServe()
	if serveErr != nil {
		return fmt.Errorf("listen and serve: %w", serveErr)
	}

	return nil
}

//nolint:ireturn
func newRepository(cfg *config.Config) (repository.Repository, error) {
	switch cfg.Storage {
	case "postgres":
		sqlDB, err := sql.Open("pgx", cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("open db: %w", err)
		}

		pingErr := sqlDB.PingContext(context.Background())
		if pingErr != nil {
			return nil, fmt.Errorf("ping db: %w", pingErr)
		}

		migrateErr := migrator.Run(sqlDB)
		if migrateErr != nil {
			return nil, fmt.Errorf("run migrations: %w", migrateErr)
		}

		return postgres.New(sqlDB), nil

	case "memory":
		return memory.New(), nil

	default:
		return nil, fmt.Errorf("%w: %q", errUnknownStorage, cfg.Storage)
	}
}
