package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"url-shortener/internal/services"
)

type Handler struct {
	Service *services.URLService
}

type shortenRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"customCode,omitempty"`
}

type shortenResponse struct {
	ShortCode   string `json:"shortCode"`
	ShortURL    string `json:"shortUrl"`
	OriginalURL string `json:"originalUrl"`
}

// POST /shorten
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.Service.Shorten(req.URL, req.CustomCode, time.Hour*24*365)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := shortenResponse{
		ShortCode:   url.ShortCode,
		ShortURL:    "http://localhost:8000/" + url.ShortCode, // later replace with BASE_URL from env
		OriginalURL: url.OriginalURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /{code}
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	url, err := h.Service.Resolve(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

// GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
