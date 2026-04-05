package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	webpush "github.com/SherClockHolmes/webpush-go"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

// NotifyHandler handles push notification dispatch.
type NotifyHandler struct {
	inventoryRepo repository.InventoryRepo
	subRepo       repository.SubscriptionRepo
}

// NewNotifyHandler creates a new notification handler.
func NewNotifyHandler(inventoryRepo repository.InventoryRepo, subRepo repository.SubscriptionRepo) *NotifyHandler {
	return &NotifyHandler{inventoryRepo: inventoryRepo, subRepo: subRepo}
}

// NotifyResult is the response body for the inventory alert endpoint.
type NotifyResult struct {
	Message  string   `json:"message"`
	Targets  []string `json:"targets"`
	Notified int      `json:"notified"`
	Failed   int      `json:"failed"`
}

// InventoryAlert sends WebPush notifications to all subscribers when blood inventory is urgent.
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

	urgentTypes := h.findUrgentBloodTypes()
	if len(urgentTypes) == 0 {
		writeJSON(w, http.StatusOK, NotifyResult{
			Notified: 0,
			Message:  "全血液型の在庫は正常です",
		})
		return
	}

	subs := h.subRepo.List()
	if len(subs) == 0 {
		writeJSON(w, http.StatusOK, NotifyResult{
			Notified: 0,
			Message:  "購読者がいません",
			Targets:  urgentTypes,
		})
		return
	}

	opts, err := loadVAPIDOptions()
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "VAPID keys not configured")
		return
	}

	payloadJSON, err := buildAlertPayload(urgentTypes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to build notification payload")
		return
	}

	sent, failed := h.dispatchNotifications(subs, payloadJSON, opts)

	writeJSON(w, http.StatusOK, NotifyResult{
		Notified: sent,
		Failed:   failed,
		Message:  "在庫逼迫通知を " + strconv.Itoa(sent) + " 件の購読者に送信しました",
		Targets:  urgentTypes,
	})
}

func (h *NotifyHandler) findUrgentBloodTypes() []string {
	inventories := h.inventoryRepo.List()
	var urgentTypes []string
	for _, inv := range inventories {
		if inv.Level == "critical" || inv.Level == "low" {
			urgentTypes = append(urgentTypes, string(inv.BloodType))
		}
	}
	return urgentTypes
}

func loadVAPIDOptions() (*webpush.Options, error) {
	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")
	vapidSubject := os.Getenv("VAPID_SUBJECT")

	if vapidPublicKey == "" || vapidPrivateKey == "" || vapidSubject == "" {
		return nil, http.ErrNotSupported
	}

	return &webpush.Options{
		Subscriber:      vapidSubject,
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             60,
	}, nil
}

func buildAlertPayload(urgentTypes []string) ([]byte, error) {
	payload := map[string]string{
		"title": "血液在庫逼迫アラート",
		"body":  "以下の血液型の在庫が逼迫しています: " + strings.Join(urgentTypes, ", "),
		"url":   "/inventory",
	}
	return json.Marshal(payload)
}

func (h *NotifyHandler) dispatchNotifications(
	subs []*model.PushSubscription,
	payloadJSON []byte,
	opts *webpush.Options,
) (sent, failed int) {
	for _, sub := range subs {
		subscription := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}

		resp, sendErr := webpush.SendNotification(payloadJSON, subscription, opts)
		if sendErr != nil {
			log.Printf("WebPush send error for %s: %v", sub.Endpoint, sendErr)
			failed++
			continue
		}
		if resp == nil {
			sent++
			continue
		}

		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close WebPush response body: %v", err)
		}

		if h.isExpiredSubscription(resp.StatusCode) {
			log.Printf("Subscription expired (status %d), removing: %s", resp.StatusCode, sub.Endpoint)
			if delErr := h.subRepo.DeleteByEndpoint(sub.Endpoint); delErr != nil {
				log.Printf("Failed to delete expired subscription: %v", delErr)
			}
			failed++
			continue
		}

		if resp.StatusCode >= http.StatusBadRequest {
			log.Printf("WebPush server returned status %d for %s", resp.StatusCode, sub.Endpoint)
			failed++
			continue
		}

		sent++
	}
	return sent, failed
}

func (h *NotifyHandler) isExpiredSubscription(statusCode int) bool {
	return statusCode == http.StatusGone || statusCode == http.StatusNotFound
}
