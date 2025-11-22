package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/models"
	"url-shortener/internal/shortener"
)

func TestHandleShorten(t *testing.T) {
	service := shortener.NewService("http://localhost:8080")
	handler := NewHandler(service)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid URL",
			requestBody:    `{"long_url":"https://www.google.com"}`,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Empty URL",
			requestBody:    `{"long_url":""}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{invalid}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleShorten(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var resp models.ShortenResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if resp.ShortURL == "" {
					t.Error("Expected non-empty short URL")
				}
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	service := shortener.NewService("http://localhost:8080")
	handler := NewHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HandleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected 'OK', got '%s'", w.Body.String())
	}
}

func TestHandleStats(t *testing.T) {
	service := shortener.NewService("http://localhost:8080")
	handler := NewHandler(service)

	// Add some URLs
	service.Shorten("https://www.google.com")
	service.Shorten("https://www.github.com")

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	handler.HandleStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp models.StatsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if resp.TotalURLsShortened != 2 {
		t.Errorf("Expected 2 URLs, got %d", resp.TotalURLsShortened)
	}

	if resp.TotalRedirects != 0 {
		t.Errorf("Expected 0 redirects, got %d", resp.TotalRedirects)
	}
}
