package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"

	"github.com/6ermvH/url-shortener/internal/repository"
	"github.com/6ermvH/url-shortener/pkg/base63"
)

const (
	shortURLLen = 10
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var (
	ErrNotFound   = errors.New("url not found")
	ErrEmptyURL   = errors.New("url is required")
	ErrInvalidURL = errors.New("invalid url")
)

type Service struct {
	repo     repository.Repository
	encoding *base63.Encoding
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo:     repo,
		encoding: base63.NewEncoding(alphabet, shortURLLen),
	}
}

func (s *Service) Shorten(ctx context.Context, req ShortenRequest) (ShortenResponse, error) {
	err := validateURL(req.URL)
	if err != nil {
		return ShortenResponse{}, err
	}

	for attempt := 0; ; attempt++ {
		short := s.generateShort(req.URL, attempt)

		err := s.repo.Save(ctx, repository.URLMapping{
			ShortURL:    short,
			OriginalURL: req.URL,
		})
		if errors.Is(err, repository.ErrConflict) {
			continue
		}

		if err != nil {
			return ShortenResponse{}, fmt.Errorf("save url mapping: %w", err)
		}

		return ShortenResponse{ShortURL: short}, nil
	}
}

func (s *Service) Resolve(ctx context.Context, short string) (ResolveResponse, error) {
	mapping, err := s.repo.GetByShort(ctx, short)
	if errors.Is(err, repository.ErrNotFound) {
		return ResolveResponse{}, ErrNotFound
	}

	if err != nil {
		return ResolveResponse{}, fmt.Errorf("get by short: %w", err)
	}

	return ResolveResponse{OriginalURL: mapping.OriginalURL}, nil
}

func (s *Service) generateShort(rawURL string, attempt int) string {
	input := rawURL
	if attempt > 0 {
		input = fmt.Sprintf("%s#%d", rawURL, attempt)
	}

	hash := sha256.Sum256([]byte(input))

	return s.encoding.EncodeToString(hash[:8])
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return ErrEmptyURL
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("%w: %q", ErrInvalidURL, rawURL)
	}

	return nil
}
