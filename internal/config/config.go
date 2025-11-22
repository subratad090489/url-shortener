package config

import (
	"fmt"
	"os"
)

// Config is the configuration for the application
type Config struct {
	// Port is the port the server listens on
	Port    string
	BaseURL string
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig loads the configuration from the
// environment variables or defaults
func LoadConfig() *Config {
	port := getenv("PORT", "8080")
	host := getenv("HOST", "localhost")

	baseURL := fmt.Sprintf("http://%s:%s", host, port)

	return &Config{
		Port:    port,
		BaseURL: baseURL,
	}
}
