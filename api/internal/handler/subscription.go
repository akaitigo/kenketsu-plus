package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

type SubscriptionHandler struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sub model.PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.repo.Create(&sub)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
