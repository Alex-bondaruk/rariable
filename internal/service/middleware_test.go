package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alex-bondaruk/rariable/internal/rarible"
)

func TestRecoveryMiddleware(t *testing.T) {
	fc := &fakeClient{panicOnGetOwnership: true}
	router := NewRouter(NewHandler(fc, testLogger()))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ownerships/ETHEREUM:0xabc:1:0xdef", nil)

	// If recovery doesn't work, this call panics and fails the test process.
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500, body=%s", rec.Code, rec.Body)
	}
}

func TestLoggingMiddlewareDoesNotAlterResponse(t *testing.T) {
	fc := &fakeClient{ownership: &rarible.Ownership{ID: "x"}}
	router := NewRouter(NewHandler(fc, testLogger()))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ownerships/ETHEREUM:0xabc:1:0xdef", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
