package config

import (
	"errors"
	"os"
)

type Config struct {
	APIKey   string
	BaseURL  string
	HTTPPort string
}

func Load() (*Config, error) {
	rariableKey := os.Getenv("RARIBLE_API_KEY")
	if rariableKey == "" {
		return nil, errors.New("RARIBLE_API_KEY is required")
	}
	rariableBaseURL := os.Getenv("RARIBLE_BASE_URL")
	if rariableBaseURL == "" {
		rariableBaseURL = "https://api.rarible.org"
	}
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		APIKey:   rariableKey,
		BaseURL:  rariableBaseURL,
		HTTPPort: port,
	}, nil
}
