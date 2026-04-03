package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

type InventoryHandler struct {
	repo *repository.InventoryRepository
}

func NewInventoryHandler(repo *repository.InventoryRepository) *InventoryHandler {
	return &InventoryHandler{repo: repo}
}

func (h *InventoryHandler) List(w http.ResponseWriter, _ *http.Request) {
	inventory := h.repo.List()
	writeJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	bloodType := model.BloodType(r.PathValue("bloodType"))

	var body struct {
		Level model.InventoryLevel `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.repo.Update(bloodType, body.Level)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
