package main

import (
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/shortener"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create the URL shortener service
	service := shortener.NewService(cfg.BaseURL)

	// Create the HTTP handler
	handler := handlers.NewHandler(service)

	// Set up HTTP routes
	http.HandleFunc("/shorten", handler.HandleShorten)
	http.HandleFunc("/stats", handler.HandleStats)
	http.HandleFunc("/", handler.HandleRedirect)
	http.HandleFunc("/health", handler.HandleHealth)

	// Start the HTTP server
	log.Printf("URL Shortener service starting on port %s", cfg.Port)
	log.Printf("Base URL: %s", cfg.BaseURL)
	log.Printf("Endpoints:")
	log.Printf("  POST %s/shorten - Shorten a URL", cfg.BaseURL)
	log.Printf("  GET  %s/{code} - Redirect to long URL", cfg.BaseURL)
	log.Printf("  GET  %s/health - Health check", cfg.BaseURL)
	log.Printf("  GET  %s/stats - Service statistics", cfg.BaseURL)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
