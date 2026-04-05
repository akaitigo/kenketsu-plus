package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

// SubscriptionHandler handles HTTP requests for push subscriptions.
type SubscriptionHandler struct {
	repo repository.SubscriptionRepo
}

// NewSubscriptionHandler creates a new handler for push subscription endpoints.
func NewSubscriptionHandler(repo repository.SubscriptionRepo) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

// Create registers a new push subscription.
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sub model.PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.repo.Create(r.Context(), &sub)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// Delete removes a push subscription by its ID.
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
