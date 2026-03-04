package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

const readHeaderTimeout = 5 * time.Second

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	repo, err := newRepository(cfg)
	if err != nil {
		return fmt.Errorf("init repository: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	svc := service.New(repo)
	h := handler.New(svc, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.Shorten)
	mux.HandleFunc("GET /{short}", h.Resolve)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger.Info("starting server", "addr", cfg.Addr, "storage", cfg.Storage)

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

		migrateErr := migrator.Run(sqlDB, cfg.MigrateVersion)
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
