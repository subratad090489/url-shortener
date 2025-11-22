package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"url-shortener/internal/models"
	"url-shortener/internal/shortener"
)

// Handler manages HTTP requests for URL shortening
type Handler struct {
	service *shortener.Service
}

// NewHandler creates a new HTTP handler for URL shortening
func NewHandler(service *shortener.Service) *Handler {
	return &Handler{service: service}
}

// respondJSON sends a JSON response
func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, message string, status int) {
	resp := models.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}
	h.respondJSON(w, resp, status)
}

// HandleShorten handles POST requests to shorten URLs
func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.LongURL == "" {
		log.Printf("Long URL is required")
		http.Error(w, "Long URL is required", http.StatusBadRequest)
		return
	}

	shortURL := h.service.Shorten(req.LongURL)

	resp := models.ShortenResponse{
		ShortURL: shortURL,
		LongURL:  req.LongURL,
	}
	h.respondJSON(w, resp, http.StatusOK)
}

// HandleRedirect handles GET requests to redirect to long URLs
func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/"):] // Extract short URL from path

	if shortURL == "" || shortURL == "shorten" ||
		shortURL == "health" || shortURL == "stats" {
		log.Printf("ShortURL Not found: %s", shortURL)
		h.respondError(w, "Not found", http.StatusNotFound)
		return
	}

	longURL, exists := h.service.GetLongURL(shortURL)
	if !exists {
		log.Printf("Short URL not found in the service: %s", shortURL)
		h.respondError(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

// HandleHealth handles health check requests
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// HandleStats handles statistics requests
func (h *Handler) HandleStats(w http.ResponseWriter, r *http.Request) {
	resp := models.StatsResponse{
		TotalURLsShortened: h.service.TotalURLShortened(),
		TotalRedirects:     h.service.TotalURLRedirects(),
	}
	h.respondJSON(w, resp, http.StatusOK)
}
