package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/handler"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func newTestDonationRouter() http.Handler {
	return handler.NewRouter(
		repository.NewCenterRepository(),
		repository.NewDonationRepository(),
		repository.NewInventoryRepository(),
		repository.NewSubscriptionRepository(),
		service.NewDonationCalculator(),
	)
}

func TestDonation_ListEmpty(t *testing.T) {
	t.Parallel()
	router := newTestDonationRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/donations", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestDonation_CreateValid(t *testing.T) {
	t.Parallel()
	router := newTestDonationRouter()

	body := `{"bloodType":"A+","donationType":"whole_400","gender":"male","donatedAt":"2026-03-01T10:00:00Z","volumeMl":400}`
	req := httptest.NewRequest(http.MethodPost, "/api/donations", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["id"] == nil {
		t.Error("expected id in response")
	}
}

func TestDonation_CreateInvalidBloodType(t *testing.T) {
	t.Parallel()
	router := newTestDonationRouter()

	body := `{"bloodType":"X+","donationType":"whole_400","gender":"male","donatedAt":"2026-03-01T10:00:00Z","volumeMl":400}`
	req := httptest.NewRequest(http.MethodPost, "/api/donations", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestDonation_NextAvailable_NoGender(t *testing.T) {
	t.Parallel()
	router := newTestDonationRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/donations/next-available", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestDonation_NextAvailable_Valid(t *testing.T) {
	t.Parallel()
	router := newTestDonationRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/donations/next-available?gender=male", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if result["canDonateToday"] != true {
		t.Error("expected canDonateToday=true with no donations")
	}
}
