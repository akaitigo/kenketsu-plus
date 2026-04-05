package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotifyInventoryAlert_NoSecret(t *testing.T) {
	router := newTestRouter()

	// NOTIFY_SECRET not set → 503
	t.Setenv("NOTIFY_SECRET", "")
	req := httptest.NewRequest(http.MethodPost, "/api/notify/inventory-alert", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestNotifyInventoryAlert_Unauthorized(t *testing.T) {
	router := newTestRouter()

	t.Setenv("NOTIFY_SECRET", "test-secret")
	req := httptest.NewRequest(http.MethodPost, "/api/notify/inventory-alert", nil)
	req.Header.Set("X-Notify-Secret", "wrong-secret")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestNotifyInventoryAlert_NoUrgentTypes(t *testing.T) {
	router := newTestRouter()

	t.Setenv("NOTIFY_SECRET", "test-secret")
	req := httptest.NewRequest(http.MethodPost, "/api/notify/inventory-alert", nil)
	req.Header.Set("X-Notify-Secret", "test-secret")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	notified, ok := result["notified"].(float64)
	if !ok || notified != 0 {
		t.Errorf("expected notified=0, got %v", result["notified"])
	}
}

func TestNotifyInventoryAlert_NoVAPID(t *testing.T) {
	// Set up inventory with critical level
	router := newTestNotifyRouter(t)

	t.Setenv("NOTIFY_SECRET", "test-secret")
	t.Setenv("VAPID_PUBLIC_KEY", "")
	t.Setenv("VAPID_PRIVATE_KEY", "")
	t.Setenv("VAPID_SUBJECT", "")

	req := httptest.NewRequest(http.MethodPost, "/api/notify/inventory-alert", nil)
	req.Header.Set("X-Notify-Secret", "test-secret")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestNotifyInventoryAlert_NoSubscribers(t *testing.T) {
	router := newTestNotifyRouterNoSubs(t)

	t.Setenv("NOTIFY_SECRET", "test-secret")
	t.Setenv("VAPID_PUBLIC_KEY", "test-key")
	t.Setenv("VAPID_PRIVATE_KEY", "test-key")
	t.Setenv("VAPID_SUBJECT", "mailto:test@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/notify/inventory-alert", nil)
	req.Header.Set("X-Notify-Secret", "test-secret")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if result["message"] != "購読者がいません" {
		t.Errorf("expected no subscribers message, got: %v", result["message"])
	}
}

// newTestNotifyRouter sets up a router with critical inventory and a subscription.
func newTestNotifyRouter(t *testing.T) http.Handler {
	t.Helper()
	router := newTestRouter()

	// Update inventory to critical
	body := `{"level":"critical"}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/inventory/A+", bytes.NewBufferString(body))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	// Add subscription
	subBody := `{"endpoint":"https://push.example.com/test","p256dh":"test-p256dh","auth":"test-auth"}`
	subReq := httptest.NewRequest(http.MethodPost, "/api/subscriptions", bytes.NewBufferString(subBody))
	subRec := httptest.NewRecorder()
	router.ServeHTTP(subRec, subReq)

	return router
}

// newTestNotifyRouterNoSubs sets up a router with critical inventory but no subscription.
func newTestNotifyRouterNoSubs(t *testing.T) http.Handler {
	t.Helper()
	router := newTestRouter()

	// Update inventory to critical
	body := `{"level":"critical"}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/inventory/A+", bytes.NewBufferString(body))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	return router
}
