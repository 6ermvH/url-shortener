package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/6ermvH/url-shortener/internal/repository"
)

//go:embed queries/get_by_short.sql
var queryGetByShort string

//go:embed queries/save.sql
var querySave string

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByShort(ctx context.Context, short string) (repository.URLMapping, error) {
	var original string

	err := r.db.QueryRowContext(ctx, queryGetByShort, short).Scan(&original)
	if errors.Is(err, sql.ErrNoRows) {
		return repository.URLMapping{}, repository.ErrNotFound
	}

	if err != nil {
		return repository.URLMapping{}, fmt.Errorf("query get by short: %w", err)
	}

	return repository.URLMapping{ShortURL: short, OriginalURL: original}, nil
}

func (r *Repository) Save(ctx context.Context, mapping repository.URLMapping) error {
	_, err := r.db.ExecContext(ctx, querySave, mapping.ShortURL, mapping.OriginalURL)
	if err != nil {
		return fmt.Errorf("query save: %w", err)
	}

	return nil
}
