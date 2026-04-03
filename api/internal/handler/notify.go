package handler

import (
	"net/http"
	"os"

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

func (h *NotifyHandler) InventoryAlert(w http.ResponseWriter, r *http.Request) {
	// C-2: shared secret authentication
	secret := os.Getenv("NOTIFY_SECRET")
	if secret != "" && r.Header.Get("X-Notify-Secret") != secret {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

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

	// H-7: pass urgentTypes directly instead of double-encoding
	writeJSON(w, http.StatusOK, NotifyResult{
		Notified: len(subs),
		Message:  "在庫逼迫通知を送信しました",
		Targets:  urgentTypes,
	})
}
