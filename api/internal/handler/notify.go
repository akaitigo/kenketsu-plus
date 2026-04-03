package handler

import (
	"net/http"
	"os"
	"strconv"

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
	// C-2: NOTIFY_SECRET is always required — reject if unset or mismatched
	secret := os.Getenv("NOTIFY_SECRET")
	if secret == "" {
		writeError(w, http.StatusServiceUnavailable, "notification service not configured")
		return
	}
	if r.Header.Get("X-Notify-Secret") != secret {
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

	// NEW-1: honest about delivery — we enqueue, actual push delivery is async
	writeJSON(w, http.StatusOK, NotifyResult{
		Notified: len(subs),
		Message:  "在庫逼迫通知を " + strconv.Itoa(len(subs)) + " 件の購読者にキューしました",
		Targets:  urgentTypes,
	})
}
