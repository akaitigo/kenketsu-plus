package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

// InventoryHandler handles HTTP requests for blood inventory.
type InventoryHandler struct {
	repo repository.InventoryRepo
}

// NewInventoryHandler creates a new handler for blood inventory endpoints.
func NewInventoryHandler(repo repository.InventoryRepo) *InventoryHandler {
	return &InventoryHandler{repo: repo}
}

// List returns all blood inventory records.
func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	inventory := h.repo.List(r.Context())
	writeJSON(w, http.StatusOK, inventory)
}

// Update changes the inventory level for a specified blood type.
func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	bloodType := model.BloodType(r.PathValue("bloodType"))

	var body struct {
		Level model.InventoryLevel `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.repo.Update(r.Context(), bloodType, body.Level)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
