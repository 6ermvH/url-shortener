package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/6ermvH/url-shortener/internal/repository"
	"github.com/6ermvH/url-shortener/internal/repository/mocks"
	"github.com/6ermvH/url-shortener/internal/service"
)

func TestShorten_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	svc := service.New(repo)

	resp, err := svc.Shorten(context.Background(), service.ShortenRequest{URL: "https://example.com"})

	require.NoError(t, err)
	require.Len(t, resp.ShortURL, 10)
}

func TestShorten_Idempotent(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	svc := service.New(repo)

	resp1, err := svc.Shorten(context.Background(), service.ShortenRequest{URL: "https://example.com"})
	require.NoError(t, err)

	resp2, err := svc.Shorten(context.Background(), service.ShortenRequest{URL: "https://example.com"})
	require.NoError(t, err)

	require.Equal(t, resp1.ShortURL, resp2.ShortURL)
}

func TestShorten_EmptyURL(t *testing.T) {
	t.Parallel()

	svc := service.New(nil)

	_, err := svc.Shorten(context.Background(), service.ShortenRequest{URL: ""})

	require.ErrorIs(t, err, service.ErrEmptyURL)
}

func TestShorten_InvalidURL(t *testing.T) {
	t.Parallel()

	svc := service.New(nil)

	_, err := svc.Shorten(context.Background(), service.ShortenRequest{URL: "not-a-url"})

	require.ErrorIs(t, err, service.ErrInvalidURL)
}

func TestShorten_RepoError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)

	repoErr := errors.New("db error")
	repo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(repoErr)

	svc := service.New(repo)

	_, err := svc.Shorten(context.Background(), service.ShortenRequest{URL: "https://example.com"})

	require.ErrorIs(t, err, repoErr)
}

func TestResolve_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().
		GetByShort(gomock.Any(), "abc123").
		Return(repository.URLMapping{ShortURL: "abc123", OriginalURL: "https://example.com"}, nil)

	svc := service.New(repo)

	resp, err := svc.Resolve(context.Background(), "abc123")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", resp.OriginalURL)
}

func TestResolve_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().
		GetByShort(gomock.Any(), "notexists").
		Return(repository.URLMapping{}, repository.ErrNotFound)

	svc := service.New(repo)

	_, err := svc.Resolve(context.Background(), "notexists")

	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestResolve_RepoError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)

	repoErr := errors.New("db error")
	repo.EXPECT().
		GetByShort(gomock.Any(), "abc123").
		Return(repository.URLMapping{}, repoErr)

	svc := service.New(repo)

	_, err := svc.Resolve(context.Background(), "abc123")

	require.ErrorIs(t, err, repoErr)
}
