package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/handler"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	router := handler.NewRouter(
		repository.NewCenterRepository(),
		repository.NewDonationRepository(),
		repository.NewInventoryRepository(),
		repository.NewSubscriptionRepository(),
		service.NewDonationCalculator(),
	)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
