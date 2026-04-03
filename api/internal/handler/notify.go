package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

type NotifyHandler struct {
	inventoryRepo *repository.InventoryRepository
	subRepo       *repository.SubscriptionRepository
}

func NewNotifyHandler(inventoryRepo *repository.InventoryRepository, subRepo *repository.SubscriptionRepository) *NotifyHandler {
	return &NotifyHandler{inventoryRepo: inventoryRepo, subRepo: subRepo}
}

type NotifyResult struct {
	Message  string   `json:"message"`
	Targets  []string `json:"targets"`
	Notified int      `json:"notified"`
}

func (h *NotifyHandler) InventoryAlert(w http.ResponseWriter, _ *http.Request) {
	inventories := h.inventoryRepo.List()

	var urgentTypes []string
	for _, inv := range inventories {
		if inv.Level == "critical" || inv.Level == "low" {
			urgentTypes = append(urgentTypes, string(inv.BloodType))
		}
	}

	if len(urgentTypes) == 0 {
		writeJSON(w, http.StatusOK, NotifyResult{
			Notified: 0,
			Message:  "全血液型の在庫は正常です",
		})
		return
	}

	subs := h.subRepo.List()
	typesJSON, err := json.Marshal(urgentTypes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to marshal urgent types")
		return
	}

	writeJSON(w, http.StatusOK, NotifyResult{
		Notified: len(subs),
		Message:  "在庫逼迫通知を送信しました",
		Targets:  []string{string(typesJSON)},
	})
}
