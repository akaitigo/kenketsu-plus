package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/handler"
)

// okHandler is a trivial handler used to detect whether RequireAdminKey passed
// the request through to the wrapped handler.
func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func doAdminRequest(t *testing.T, key string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/centers", nil)
	if key != "" {
		req.Header.Set("X-Admin-Key", key)
	}
	rec := httptest.NewRecorder()
	handler.RequireAdminKey(okHandler)(rec, req)
	return rec
}

// TestRequireAdminKey_Unset verifies the fail-closed behavior (#20): when
// ADMIN_API_KEY is unset the endpoint is rejected with 503 instead of bypassed.
func TestRequireAdminKey_Unset(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "")

	rec := doAdminRequest(t, "anything")

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 when ADMIN_API_KEY unset, got %d", rec.Code)
	}
}

func TestRequireAdminKey_CorrectKey(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "secret-key")

	rec := doAdminRequest(t, "secret-key")

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 with correct key, got %d", rec.Code)
	}
}

func TestRequireAdminKey_WrongKey(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "secret-key")

	rec := doAdminRequest(t, "wrong-key")

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 with wrong key, got %d", rec.Code)
	}
}

func TestRequireAdminKey_MissingHeader(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "secret-key")

	rec := doAdminRequest(t, "")

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 with missing header, got %d", rec.Code)
	}
}
