package config

import "testing"

func TestLoad(t *testing.T) {
	t.Run("return error when api key is missing", func(t *testing.T) {
		t.Setenv("RARIBLE_API_KEY", "")
		t.Setenv("RARIBLE_BASE_URL", "")
		t.Setenv("HTTP_PORT", "")
		cfg, err := Load()
		if err == nil {
			t.Fatal("expected error when RARIBLE_API_KEY is missing, got nil")
		}
		if cfg != nil {
			t.Errorf("expected nil config, got %+v", cfg)
		}
	})

	t.Run("applies defaults when only api key is set", func(t *testing.T) {
		t.Setenv("RARIBLE_API_KEY", "test-key")
		t.Setenv("RARIBLE_BASE_URL", "")
		t.Setenv("HTTP_PORT", "")
		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected config, got nil")
		}
		if cfg.BaseURL != "https://api.rarible.org" {
			t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://api.rarible.org")
		}
		if cfg.HTTPPort != "8080" {
			t.Errorf("HTTPPort = %q, want %q", cfg.HTTPPort, "8080")
		}
		if cfg.APIKey != "test-key" {
			t.Errorf("APIKey = %q, want %q", cfg.APIKey, "test-key")
		}
	})

	t.Run("uses provided values", func(t *testing.T) {
		t.Setenv("RARIBLE_API_KEY", "test-new-key")
		t.Setenv("RARIBLE_BASE_URL", "https://testnet-api.rarible.org")
		t.Setenv("HTTP_PORT", "8081")
		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected config, got nil")
		}
		if cfg.BaseURL != "https://testnet-api.rarible.org" {
			t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://testnet-api.rarible.org")
		}
		if cfg.HTTPPort != "8081" {
			t.Errorf("HTTPPort = %q, want %q", cfg.HTTPPort, "8081")
		}
		if cfg.APIKey != "test-new-key" {
			t.Errorf("APIKey = %q, want %q", cfg.APIKey, "test-new-key")
		}
	})
}
