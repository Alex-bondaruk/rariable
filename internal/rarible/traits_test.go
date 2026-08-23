package rarible

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryTraitsWithRarity(t *testing.T) {
	validReq := TraitsRarityRequest{
		CollectionID: "ETHEREUM:0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
		Properties:   []TraitProperty{{Key: "Hat", Value: "Halo"}},
	}

	t.Run("decodes a successful response", func(t *testing.T) {
		srv := serveFixture(t, http.StatusOK, "traits_ok.json", nil)

		c, _ := New(srv.URL, "test-key")
		res, err := c.QueryTraitsWithRarity(context.Background(), validReq)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(res.Traits) != 2 {
			t.Fatalf("len(Traits) = %d, want 2", len(res.Traits))
		}
		want := TraitRarity{Key: "Hat", Value: "Halo", Rarity: "3.2406481"}
		if res.Traits[0] != want {
			t.Errorf("Traits[0] = %+v, want %+v", res.Traits[0], want)
		}
	})

	t.Run("posts the request body as json", func(t *testing.T) {
		var gotBody []byte
		var gotReq *http.Request
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotReq = r.Clone(r.Context())
			gotBody, _ = io.ReadAll(r.Body)
			w.Write([]byte(`{"traits":[]}`))
		}))
		t.Cleanup(srv.Close)

		c, _ := New(srv.URL, "test-key")
		if _, err := c.QueryTraitsWithRarity(context.Background(), validReq); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if gotReq.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", gotReq.Method)
		}
		if ct := gotReq.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}

		var sent TraitsRarityRequest
		if err := json.Unmarshal(gotBody, &sent); err != nil {
			t.Fatalf("body is not valid json: %v", err)
		}
		if sent.CollectionID != validReq.CollectionID {
			t.Errorf("collectionId = %q, want %q", sent.CollectionID, validReq.CollectionID)
		}
		if len(sent.Properties) != 1 || sent.Properties[0] != validReq.Properties[0] {
			t.Errorf("properties = %+v, want %+v", sent.Properties, validReq.Properties)
		}
	})

	t.Run("maps a non-2xx response to APIError", func(t *testing.T) {
		srv := serveFixture(t, http.StatusBadRequest, "ownership_badrequest.json", nil)

		c, _ := New(srv.URL, "test-key")
		_, err := c.QueryTraitsWithRarity(context.Background(), validReq)

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T: %v", err, err)
		}
		if apiErr.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
		}
	})

	t.Run("rejects invalid input without calling the api", func(t *testing.T) {
		called := false
		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			called = true
		}))
		t.Cleanup(srv.Close)

		c, _ := New(srv.URL, "test-key")

		cases := []TraitsRarityRequest{
			{Properties: []TraitProperty{{Key: "Hat", Value: "Halo"}}},
			{CollectionID: "ETHEREUM:0xabc"},
		}
		for _, req := range cases {
			if _, err := c.QueryTraitsWithRarity(context.Background(), req); err == nil {
				t.Errorf("expected error for %+v, got nil", req)
			}
		}
		if called {
			t.Error("the api was called for invalid input")
		}
	})
}
