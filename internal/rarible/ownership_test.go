package rarible

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const testOwnershipID = "ETHEREUM:0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d:664:0x4459084da2d3a774c436f2e75f2e3fe9335dc5de"

// serveFixture starts a server that replies with status and the given fixture file.
func serveFixture(t *testing.T, status int, fixture string, gotReq **http.Request) *httptest.Server {
	t.Helper()

	body, err := os.ReadFile("testdata/" + fixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if gotReq != nil {
			*gotReq = r.Clone(r.Context())
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	return srv
}

func TestGetOwnershipByID(t *testing.T) {
	t.Run("decodes a successful response", func(t *testing.T) {
		var got *http.Request
		srv := serveFixture(t, http.StatusOK, "ownership_ok.json", &got)

		c, err := New(srv.URL, "test-key")
		if err != nil {
			t.Fatalf("New: %v", err)
		}

		o, err := c.GetOwnershipByID(context.Background(), testOwnershipID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if o == nil {
			t.Fatal("expected ownership, got nil")
		}

		if o.ID != testOwnershipID {
			t.Errorf("ID = %q, want %q", o.ID, testOwnershipID)
		}
		if o.TokenID != "664" {
			t.Errorf("TokenID = %q, want %q", o.TokenID, "664")
		}
		if o.Value != "1" {
			t.Errorf("Value = %q, want %q", o.Value, "1")
		}
		if o.CreatedAt.IsZero() {
			t.Error("CreatedAt is zero, want a parsed timestamp")
		}
		if o.BestSellOrder == nil {
			t.Fatal("expected BestSellOrder, got nil")
		}
		if o.BestSellOrder.MakePrice != "7.849999" {
			t.Errorf("MakePrice = %q, want %q", o.BestSellOrder.MakePrice, "7.849999")
		}
	})

	t.Run("sends the api key and the escaped id", func(t *testing.T) {
		var got *http.Request
		srv := serveFixture(t, http.StatusOK, "ownership_ok.json", &got)

		c, _ := New(srv.URL, "test-key")
		if _, err := c.GetOwnershipByID(context.Background(), testOwnershipID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if key := got.Header.Get("X-API-KEY"); key != "test-key" {
			t.Errorf("X-API-KEY = %q, want %q", key, "test-key")
		}
		wantPath := "/v0.1/ownerships/" + url.PathEscape(testOwnershipID)
		if got.URL.EscapedPath() != wantPath {
			t.Errorf("path = %q, want %q", got.URL.EscapedPath(), wantPath)
		}
	})

	t.Run("keeps a question mark out of the query string", func(t *testing.T) {
		var got *http.Request
		srv := serveFixture(t, http.StatusOK, "ownership_ok.json", &got)

		c, _ := New(srv.URL, "test-key")
		if _, err := c.GetOwnershipByID(context.Background(), "ETHEREUM:0xabc:1:0xdef?foo=bar"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.URL.RawQuery != "" {
			t.Errorf("RawQuery = %q, want empty", got.URL.RawQuery)
		}
	})

	t.Run("maps errors to APIError", func(t *testing.T) {
		cases := []struct {
			name     string
			status   int
			fixture  string
			wantCode string
		}{
			{"not found", http.StatusNotFound, "ownership_notfound.json", "NOT_FOUND"},
			{"bad request", http.StatusBadRequest, "ownership_badrequest.json", "BAD_REQUEST"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				srv := serveFixture(t, tc.status, tc.fixture, nil)

				c, _ := New(srv.URL, "test-key")
				o, err := c.GetOwnershipByID(context.Background(), testOwnershipID)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if o != nil {
					t.Errorf("expected nil ownership, got %+v", o)
				}

				var apiErr *APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected *APIError, got %T: %v", err, err)
				}
				if apiErr.StatusCode != tc.status {
					t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tc.status)
				}
				if apiErr.Code != tc.wantCode {
					t.Errorf("Code = %q, want %q", apiErr.Code, tc.wantCode)
				}
				if apiErr.Message == "" {
					t.Error("Message is empty, want the message from Rarible")
				}
			})
		}
	})

	t.Run("rejects an empty id without calling the api", func(t *testing.T) {
		called := false
		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			called = true
		}))
		t.Cleanup(srv.Close)

		c, _ := New(srv.URL, "test-key")
		if _, err := c.GetOwnershipByID(context.Background(), ""); err == nil {
			t.Fatal("expected error, got nil")
		}
		if called {
			t.Error("the api was called for an empty id")
		}
	})
}
