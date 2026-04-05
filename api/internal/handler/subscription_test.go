package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/handler"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func newTestSubRouter() http.Handler {
	return handler.NewRouter(
		repository.NewCenterRepository(),
		repository.NewDonationRepository(),
		repository.NewInventoryRepository(),
		repository.NewSubscriptionRepository(),
		service.NewDonationCalculator(),
	)
}

func TestSubscription_Create(t *testing.T) {
	t.Parallel()
	router := newTestSubRouter()

	body := `{"endpoint":"https://push.example.com/sub1","p256dh":"test-key","auth":"test-auth"}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/subscriptions", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if result["id"] == nil {
		t.Error("expected id in response")
	}
}

func TestSubscription_CreateInvalid(t *testing.T) {
	t.Parallel()
	router := newTestSubRouter()

	body := `{"endpoint":"","p256dh":"key","auth":"auth"}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/subscriptions", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestSubscription_Delete(t *testing.T) {
	t.Parallel()
	router := newTestSubRouter()

	// Create first
	body := `{"endpoint":"https://push.example.com/sub2","p256dh":"key","auth":"auth"}`
	createReq := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/subscriptions", bytes.NewBufferString(body))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var created map[string]interface{}
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	id, ok := created["id"].(string)
	if !ok {
		t.Fatal("id not a string")
	}

	// Delete
	delReq := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/subscriptions/"+id, nil)
	delRec := httptest.NewRecorder()
	router.ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", delRec.Code)
	}
}

func TestSubscription_DeleteNotFound(t *testing.T) {
	t.Parallel()
	router := newTestSubRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/subscriptions/nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
