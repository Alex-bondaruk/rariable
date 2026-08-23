package rarible

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		apiKey  string
		wantErr bool
	}{
		{"valid arguments", "https://api.rarible.org", "test-key", false},
		{"empty api key", "https://api.rarible.org", "", true},
		{"empty base url", "", "test-key", true},
		{"base url without scheme", "mainnet", "test-key", true},
		{"malformed base url", "http://[::1", "test-key", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.baseURL, tt.apiKey)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for baseURL=%q apiKey=%q, got nil", tt.baseURL, tt.apiKey)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("expected client, got nil")
			}
		})
	}
}
