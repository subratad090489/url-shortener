package shortener

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	service := NewService("http://localhost:8080")

	tests := []struct {
		name    string
		longURL string
	}{
		{
			name:    "Google URL",
			longURL: "https://www.google.com",
		},
		{
			name:    "GitHub URL",
			longURL: "https://github.com/example/repo",
		},
		{
			name:    "Long URL",
			longURL: "https://www.example.com/very/long/path/with/many/segments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortCode1 := service.generateShortUrl(tt.longURL)
			shortCode2 := service.generateShortUrl(tt.longURL)

			// Test idempotency
			if shortCode1 != shortCode2 {
				t.Errorf("Short codes should be identical for same URL. Got %s and %s", shortCode1, shortCode2)
			}

			// Test length
			if len(shortCode1) != shortURLLength {
				t.Errorf("Expected short code length %d, got %d", shortURLLength, len(shortCode1))
			}

			// Verify Base62 characters
			for _, c := range shortCode1 {
				valid := false
				for _, validChar := range base62Chars {
					if c == validChar {
						valid = true
						break
					}
				}
				if !valid {
					t.Errorf("Invalid Base62 character in short code: %c", c)
				}
			}
		})
	}
}

func TestShortenIdempotency(t *testing.T) {
	service := NewService("http://localhost:8080")
	longURL := "https://www.example.com/test"

	shortURL1 := service.Shorten(longURL)
	shortURL2 := service.Shorten(longURL)

	if shortURL1 != shortURL2 {
		t.Errorf("Shorten should be idempotent. Got %s and %s", shortURL1, shortURL2)
	}
}

func TestGetLongURL(t *testing.T) {
	service := NewService("http://localhost:8080")
	longURL := "https://www.example.com/test"

	shortURL := service.Shorten(longURL)
	// Extract short code from full URL
	shortCode := shortURL[len("http://localhost:8080/"):]

	retrievedURL, exists := service.GetLongURL(shortCode)

	if !exists {
		t.Error("Short code should exist")
	}

	if retrievedURL != longURL {
		t.Errorf("Expected %s, got %s", longURL, retrievedURL)
	}
}

func TestCollisionHandling(t *testing.T) {
	service := NewService("http://localhost:8080")

	// Create many URLs to test collision handling
	urls := make([]string, 100)
	shortCodes := make(map[string]bool)

	for i := 0; i < 100; i++ {
		urls[i] = "https://www.example.com/test" + string(rune(i))
		shortURL := service.Shorten(urls[i])
		shortCode := shortURL[len("http://localhost:8080/"):]

		// Check for uniqueness
		if shortCodes[shortCode] {
			t.Errorf("Duplicate short code generated: %s", shortCode)
		}
		shortCodes[shortCode] = true

		// Verify retrieval
		retrievedURL, exists := service.GetLongURL(shortCode)
		if !exists {
			t.Errorf("Short code %s should exist", shortCode)
		}
		if retrievedURL != urls[i] {
			t.Errorf("URL mismatch for %s: expected %s, got %s", shortCode, urls[i], retrievedURL)
		}
	}
}
