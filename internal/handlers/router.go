package handlers

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("POST /shorten", h.ShortenURL)
	mux.HandleFunc("GET /{code}", h.Redirect)
	mux.HandleFunc("GET /health", h.Health)

	return mux
}
