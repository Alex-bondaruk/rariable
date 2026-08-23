package rarible

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

// Client talks to the Rarible Protocol API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New returns a Client for baseURL, authenticating with apiKey.
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

// do performs a request against the Rarible API.
// in is sent as a JSON body when not nil; out must be a pointer to decode a 2xx body into.
func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		raw, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("rarible: encode request body: %w", err)
		}
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("rarible: build request: %w", err)
	}

	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("rarible: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return apiErrorFrom(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("rarible: decode response: %w", err)
	}

	return nil
}
