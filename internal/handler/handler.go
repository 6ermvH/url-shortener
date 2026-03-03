package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/6ermvH/url-shortener/internal/service"
)

type Service interface {
	Shorten(ctx context.Context, input service.ShortenInput) (service.ShortenResult, error)
	Resolve(ctx context.Context, short string) (service.ResolveResult, error)
}

type Handler struct {
	svc Service
	log *slog.Logger
}

func New(svc Service, log *slog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.writeError(r.Context(), w, http.StatusBadRequest, "invalid request body")

		return
	}

	result, err := h.svc.Shorten(r.Context(), service.ShortenInput{URL: req.URL})
	if errors.Is(err, service.ErrEmptyURL) || errors.Is(err, service.ErrInvalidURL) {
		h.writeError(r.Context(), w, http.StatusBadRequest, err.Error())

		return
	}

	if err != nil {
		h.log.ErrorContext(r.Context(), "shorten url", "err", err)
		h.writeError(r.Context(), w, http.StatusInternalServerError, "internal server error")

		return
	}

	h.writeJSON(r.Context(), w, http.StatusCreated, shortenResponse{ShortURL: result.ShortURL})
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("short")

	result, err := h.svc.Resolve(r.Context(), short)
	if errors.Is(err, service.ErrNotFound) {
		h.writeError(r.Context(), w, http.StatusNotFound, "url not found")

		return
	}

	if err != nil {
		h.log.ErrorContext(r.Context(), "resolve url", "err", err)
		h.writeError(r.Context(), w, http.StatusInternalServerError, "internal server error")

		return
	}

	h.writeJSON(r.Context(), w, http.StatusOK, resolveResponse{OriginalURL: result.OriginalURL})
}

func (h *Handler) writeJSON(ctx context.Context, w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.ErrorContext(ctx, "encode response", "err", err)
	}
}

func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	h.writeJSON(ctx, w, status, errorResponse{Error: msg})
}
