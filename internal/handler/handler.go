package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/6ermvH/url-shortener/internal/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req service.ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	resp, err := h.svc.Shorten(r.Context(), req)
	if errors.Is(err, service.ErrEmptyURL) || errors.Is(err, service.ErrInvalidURL) {
		writeError(w, http.StatusBadRequest, err.Error())

		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("short")

	resp, err := h.svc.Resolve(r.Context(), short)
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, "url not found")

		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encodeErr := json.NewEncoder(w).Encode(v)
	if encodeErr != nil {
		log.Printf("writeJSON: %v", encodeErr)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
