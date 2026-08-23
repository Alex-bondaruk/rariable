package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alex-bondaruk/rariable/internal/rarible"

	"go.uber.org/zap"
)

func testLogger() *zap.Logger {
	return zap.NewNop()
}

// fakeClient is a test double for raribleClient: no HTTP, just canned results.
type fakeClient struct {
	ownership    *rarible.Ownership
	ownershipErr error
	gotID        string

	traits    *rarible.TraitsRarityResponse
	traitsErr error
	gotReq    rarible.TraitsRarityRequest

	panicOnGetOwnership bool
}

func (f *fakeClient) GetOwnershipByID(ctx context.Context, id string) (*rarible.Ownership, error) {
	if f.panicOnGetOwnership {
		panic("boom")
	}
	f.gotID = id
	return f.ownership, f.ownershipErr
}

func (f *fakeClient) QueryTraitsWithRarity(ctx context.Context, req rarible.TraitsRarityRequest) (*rarible.TraitsRarityResponse, error) {
	f.gotReq = req
	return f.traits, f.traitsErr
}

func TestGetOwnership(t *testing.T) {
	t.Run("returns the ownership as json", func(t *testing.T) {
		fc := &fakeClient{ownership: &rarible.Ownership{
			ID:      "ETHEREUM:0xabc:1:0xdef",
			TokenID: "1",
			Value:   "1",
		}}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/ETHEREUM:0xabc:1:0xdef", nil)
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body)
		}
		if fc.gotID != "ETHEREUM:0xabc:1:0xdef" {
			t.Errorf("client called with id=%q, want the full path value", fc.gotID)
		}

		var got rarible.Ownership
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("response is not valid json: %v", err)
		}
		if got.TokenID != "1" {
			t.Errorf("TokenID = %q, want %q", got.TokenID, "1")
		}
	})

	t.Run("maps a 404 from rarible to 404", func(t *testing.T) {
		fc := &fakeClient{ownershipErr: &rarible.APIError{StatusCode: http.StatusNotFound, Code: "NOT_FOUND", Message: "not found"}}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/ETHEREUM:0xabc:1:0xdef", nil)
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", rec.Code)
		}
	})

	t.Run("maps an unauthorized rarible response to 502", func(t *testing.T) {
		fc := &fakeClient{ownershipErr: &rarible.APIError{StatusCode: http.StatusForbidden, Code: "UNAUTHORIZED", Message: "bad key"}}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/ETHEREUM:0xabc:1:0xdef", nil)
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadGateway {
			t.Errorf("status = %d, want 502", rec.Code)
		}
	})

	t.Run("maps a network error to 502", func(t *testing.T) {
		fc := &fakeClient{ownershipErr: context.DeadlineExceeded}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/ETHEREUM:0xabc:1:0xdef", nil)
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadGateway {
			t.Errorf("status = %d, want 502", rec.Code)
		}
	})
}

func TestTraitsRarity(t *testing.T) {
	validBody := `{"collectionId":"ETHEREUM:0xabc","properties":[{"key":"Hat","value":"Halo"}]}`

	t.Run("returns traits as json", func(t *testing.T) {
		fc := &fakeClient{traits: &rarible.TraitsRarityResponse{
			Traits: []rarible.TraitRarity{{Key: "Hat", Value: "Halo", Rarity: "3.24"}},
		}}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", strings.NewReader(validBody))
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body)
		}
		if fc.gotReq.CollectionID != "ETHEREUM:0xabc" {
			t.Errorf("collectionId passed to client = %q", fc.gotReq.CollectionID)
		}

		var got rarible.TraitsRarityResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("response is not valid json: %v", err)
		}
		if len(got.Traits) != 1 || got.Traits[0].Rarity != "3.24" {
			t.Errorf("Traits = %+v", got.Traits)
		}
	})

	t.Run("rejects invalid json body", func(t *testing.T) {
		fc := &fakeClient{}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", strings.NewReader("not json"))
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})

	t.Run("maps a bad request from rarible to 400", func(t *testing.T) {
		fc := &fakeClient{traitsErr: &rarible.APIError{StatusCode: http.StatusBadRequest, Code: "BAD_REQUEST", Message: "bad collection"}}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", strings.NewReader(validBody))
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})

	t.Run("maps a rate limit from rarible to 503", func(t *testing.T) {
		fc := &fakeClient{traitsErr: &rarible.APIError{StatusCode: http.StatusTooManyRequests, Code: "TOO_MANY_REQUESTS"}}
		router := NewRouter(NewHandler(fc, testLogger()))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", strings.NewReader(validBody))
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want 503", rec.Code)
		}
	})
}

func TestHealthz(t *testing.T) {
	router := NewRouter(NewHandler(&fakeClient{}, testLogger()))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
