package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

func TestInventoryList(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/inventory", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var inventories []*model.BloodInventory
	if err := json.NewDecoder(rec.Body).Decode(&inventories); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(inventories) != 8 {
		t.Errorf("expected 8 blood types, got %d", len(inventories))
	}

	bloodTypes := make(map[model.BloodType]bool)
	for _, inv := range inventories {
		bloodTypes[inv.BloodType] = true
		if inv.Level != model.InventoryLevelNormal {
			t.Errorf("expected initial level to be normal, got %s for %s", inv.Level, inv.BloodType)
		}
	}

	for _, bt := range model.ValidBloodTypes {
		if !bloodTypes[bt] {
			t.Errorf("missing blood type %s in response", bt)
		}
	}
}

func TestInventoryUpdate(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	body, err := json.Marshal(map[string]string{"level": "critical"})
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/inventory/A+", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", rec.Code, rec.Body.String())
	}

	var updated model.BloodInventory
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.BloodType != model.BloodTypeAPos {
		t.Errorf("expected blood type A+, got %s", updated.BloodType)
	}
	if updated.Level != model.InventoryLevelCritical {
		t.Errorf("expected level critical, got %s", updated.Level)
	}
}

func TestInventoryUpdate_InvalidBloodType(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	body, err := json.Marshal(map[string]string{"level": "normal"})
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/inventory/X+", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestInventoryUpdate_InvalidLevel(t *testing.T) {
	t.Parallel()

	router := newTestRouter()

	body, err := json.Marshal(map[string]string{"level": "unknown"})
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/inventory/A+", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
