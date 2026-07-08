package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/handler"
	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func newTestRouter() http.Handler {
	return handler.NewRouter(
		repository.NewCenterRepository(),
		repository.NewDonationRepository(),
		repository.NewInventoryRepository(),
		repository.NewSubscriptionRepository(),
		service.NewDonationCalculator(),
	)
}

func createCenter(t *testing.T, router http.Handler, center model.DonationCenter) model.DonationCenter {
	t.Helper()

	body, err := json.Marshal(center)
	if err != nil {
		t.Fatalf("failed to marshal center: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/centers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAdminKey(req)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var created model.DonationCenter
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return created
}

func TestListCenters_Empty(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var centers []model.DonationCenter
	if err := json.NewDecoder(rec.Body).Decode(&centers); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(centers) != 0 {
		t.Errorf("expected 0 centers, got %d", len(centers))
	}
}

func TestListCenters_WithData(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	createCenter(t, router, model.DonationCenter{
		Name:           "Shibuya Center",
		Address:        "Shibuya, Tokyo",
		Lat:            35.6580,
		Lng:            139.7016,
		Capacity:       50,
		AvailableSlots: 10,
	})

	createCenter(t, router, model.DonationCenter{
		Name:           "Shinjuku Center",
		Address:        "Shinjuku, Tokyo",
		Lat:            35.6938,
		Lng:            139.7034,
		Capacity:       30,
		AvailableSlots: 5,
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var centers []model.DonationCenter
	if err := json.NewDecoder(rec.Body).Decode(&centers); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(centers) != 2 {
		t.Errorf("expected 2 centers, got %d", len(centers))
	}
}

func TestCreateCenter_Success(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	input := model.DonationCenter{
		Name:           "Ikebukuro Center",
		Address:        "Ikebukuro, Tokyo",
		Lat:            35.7295,
		Lng:            139.7109,
		Capacity:       40,
		AvailableSlots: 15,
	}

	created := createCenter(t, router, input)

	if created.ID == "" {
		t.Error("expected ID to be set")
	}
	if created.Name != input.Name {
		t.Errorf("expected name %q, got %q", input.Name, created.Name)
	}
	if created.Address != input.Address {
		t.Errorf("expected address %q, got %q", input.Address, created.Address)
	}
	if created.Status != model.CenterStatusOpen {
		t.Errorf("expected status %q, got %q", model.CenterStatusOpen, created.Status)
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected createdAt to be set")
	}
}

func TestCreateCenter_ValidationError_MissingName(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	input := model.DonationCenter{
		Address:        "Somewhere",
		Lat:            35.0,
		Lng:            139.0,
		Capacity:       10,
		AvailableSlots: 5,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/centers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAdminKey(req)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateCenter_ValidationError_InvalidLat(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	input := model.DonationCenter{
		Name:           "Bad Center",
		Address:        "Nowhere",
		Lat:            91.0,
		Lng:            139.0,
		Capacity:       10,
		AvailableSlots: 5,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/centers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAdminKey(req)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateCenter_InvalidJSON(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/centers", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	setAdminKey(req)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGetCenterByID_Success(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	created := createCenter(t, router, model.DonationCenter{
		Name:           "Tokyo Center",
		Address:        "Chiyoda, Tokyo",
		Lat:            35.6812,
		Lng:            139.7671,
		Capacity:       60,
		AvailableSlots: 20,
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers/"+created.ID, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var fetched model.DonationCenter
	if err := json.NewDecoder(rec.Body).Decode(&fetched); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected ID %q, got %q", created.ID, fetched.ID)
	}
	if fetched.Name != created.Name {
		t.Errorf("expected name %q, got %q", created.Name, fetched.Name)
	}
}

func TestGetCenterByID_NotFound(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers/nonexistent-id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestListCenters_DistanceFilter(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	// Shibuya (near the query point)
	createCenter(t, router, model.DonationCenter{
		Name:           "Shibuya Center",
		Address:        "Shibuya, Tokyo",
		Lat:            35.6580,
		Lng:            139.7016,
		Capacity:       50,
		AvailableSlots: 10,
	})

	// Osaka (far away)
	createCenter(t, router, model.DonationCenter{
		Name:           "Osaka Center",
		Address:        "Osaka",
		Lat:            34.6937,
		Lng:            135.5023,
		Capacity:       30,
		AvailableSlots: 5,
	})

	// Query near Shibuya with 5km radius
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers?lat=35.66&lng=139.70&radius=5", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var centers []model.DonationCenter
	if err := json.NewDecoder(rec.Body).Decode(&centers); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(centers) != 1 {
		t.Fatalf("expected 1 center within radius, got %d", len(centers))
	}

	if centers[0].Name != "Shibuya Center" {
		t.Errorf("expected Shibuya Center, got %q", centers[0].Name)
	}
}

func TestListCenters_DistanceFilter_InvalidLat(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers?lat=abc&lng=139.70", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestListCenters_DistanceFilter_InvalidLng(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers?lat=35.66&lng=abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestListCenters_DistanceFilter_InvalidRadius(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers?lat=35.66&lng=139.70&radius=abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestListCenters_DistanceFilter_DefaultRadius(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	// Shinjuku (about 4km from query point)
	createCenter(t, router, model.DonationCenter{
		Name:           "Shinjuku Center",
		Address:        "Shinjuku, Tokyo",
		Lat:            35.6938,
		Lng:            139.7034,
		Capacity:       30,
		AvailableSlots: 5,
	})

	// Query without radius param (default = 10km)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/centers?lat=35.66&lng=139.70", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var centers []model.DonationCenter
	if err := json.NewDecoder(rec.Body).Decode(&centers); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(centers) != 1 {
		t.Errorf("expected 1 center within default radius, got %d", len(centers))
	}
}
