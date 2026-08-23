package rarible

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string) (*Client, error) {
	baseURL = strings.TrimRight(baseURL, "/")

	if apiKey == "" {
		return nil, errors.New("rarible: apiKey is required")
	}
	if baseURL == "" {
		return nil, errors.New("rarible: baseURL is required")
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("rarible: invalid baseURL: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("rarible: baseURL must be an absolute URL, got %q", baseURL)
	}

	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}, nil
}
