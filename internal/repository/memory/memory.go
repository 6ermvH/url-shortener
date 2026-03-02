package memory

import (
	"context"
	"sync"

	"github.com/6ermvH/url-shortener/internal/repository"
)

type Repository struct {
	mu       sync.RWMutex
	byShort  map[string]string
}

func New() *Repository {
	return &Repository{
		byShort: make(map[string]string),
	}
}

func (r *Repository) GetByShort(_ context.Context, short string) (repository.URLMapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	original, ok := r.byShort[string(short)]
	if !ok {
		return repository.URLMapping{}, repository.ErrNotFound
	}

	return repository.URLMapping{ShortURL: short, OriginalURL: original}, nil
}

func (r *Repository) Save(_ context.Context, m repository.URLMapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byShort[m.ShortURL] = m.OriginalURL

	return nil
}
