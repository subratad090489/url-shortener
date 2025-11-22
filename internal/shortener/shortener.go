package shortener

import (
	"crypto/md5"
	"math/big"
	"sync"
)

// Implementation Details for creating a shortURL from a longURL
// 1. Create a hash of the long URL using MD5
// 2. Take the first 8 bytes of the hash (64 bits)
// 3. Encode the 64 bits to Base62
// 4. Take the first 7 characters of the Base62 encoded string
// 5. If there is a collision, append a counter to the short URL
//    and increment the counter until there is no collision
// 6. Return the short URL

// Service handles URL shortening logic
type Service struct {
	urlMap            map[string]string // shortURL -> longURL
	reverseMap        map[string]string // longURL -> shortURL
	mu                sync.RWMutex
	baseURL           string
	totalURLShortened int // Total URLs shortened
	totalURLRedirects int // Total URL redirects
}

const (
	// Base62 characters for encoding
	base62Chars    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortURLLength = 7 // Length of the generated short URL
)

// NewService creates a new URL shortener service
func NewService(baseURL string) *Service {
	return &Service{
		urlMap:     make(map[string]string),
		reverseMap: make(map[string]string),
		baseURL:    baseURL,
	}
}

// encodeBase62 converts a byte slice to Base62 encoding
func encodeBase62(data []byte, length int) string {
	// Convert bytes to a big integer
	num := new(big.Int).SetBytes(data)

	// Base for Base62
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	result := make([]byte, 0, length)

	// Convert to Base62
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		result = append(result, base62Chars[mod.Int64()])
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// Pad with '0' if necessary to reach desired length
	for len(result) < length {
		result = append([]byte{'0'}, result...)
	}

	// Take only the first 'length' characters
	if len(result) > length {
		result = result[:length]
	}

	return string(result)
}

// generateShortUrl creates a short URL from a long URL
// using MD5 hashing and Base62 encoding
func (s *Service) generateShortUrl(longURL string) string {
	// Create a hash of the long URL
	hash := md5.Sum([]byte(longURL))

	// Take first 8 bytes of the hash (64 bits)
	prefix := hash[:8]

	// Encode the prefix to Base62
	shortURL := encodeBase62(prefix, shortURLLength)

	return shortURL
}

// Shorten creates or retrieves a short URL for the given long URL
func (s *Service) Shorten(longURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if this long URL has already been shortened (idempotent)
	if shortURL, exists := s.reverseMap[longURL]; exists {
		return s.baseURL + "/" + shortURL
	}

	// Generate a new short URL
	shortURL := s.generateShortUrl(longURL)

	// Handle collision (very rare with MD5 + Base62)
	originalShortURL := shortURL
	counter := 1
	for {
		if _, exists := s.urlMap[shortURL]; !exists {
			break
		}
		// If collision occurs, append counter to the short URL
		shortURL = originalShortURL + string(base62Chars[counter%62])
		counter++
	}

	s.urlMap[shortURL] = longURL
	s.reverseMap[longURL] = shortURL

	s.totalURLShortened++

	return s.baseURL + "/" + shortURL
}

// GetLongURL retrieves the original long URL for a given short URL
func (s *Service) GetLongURL(shortURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	longURL, exists := s.urlMap[shortURL]
	if exists {
		s.totalURLRedirects++
	}
	return longURL, exists
}

// TotalURLShortened returns the total number of URLs shortened
func (s *Service) TotalURLShortened() int {
	return s.totalURLShortened
}

// TotalRedirects returns the total number of redirects (not implemented, placeholder)
func (s *Service) TotalURLRedirects() int {
	return s.totalURLRedirects
}
