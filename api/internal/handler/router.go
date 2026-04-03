package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func NewRouter(
	centerRepo *repository.CenterRepository,
	donationRepo *repository.DonationRepository,
	inventoryRepo *repository.InventoryRepository,
	subRepo *repository.SubscriptionRepository,
	calculator *service.DonationCalculator,
) http.Handler {
	mux := http.NewServeMux()

	centerH := NewCenterHandler(centerRepo)
	donationH := NewDonationHandler(donationRepo, calculator)
	inventoryH := NewInventoryHandler(inventoryRepo)
	subH := NewSubscriptionHandler(subRepo)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /api/centers", centerH.List)
	mux.HandleFunc("GET /api/centers/{id}", centerH.GetByID)
	mux.HandleFunc("POST /api/centers", centerH.Create)

	mux.HandleFunc("GET /api/donations", donationH.List)
	mux.HandleFunc("POST /api/donations", donationH.Create)
	mux.HandleFunc("GET /api/donations/next-available", donationH.NextAvailable)

	mux.HandleFunc("GET /api/inventory", inventoryH.List)
	mux.HandleFunc("PUT /api/inventory/{bloodType}", inventoryH.Update)

	mux.HandleFunc("POST /api/subscriptions", subH.Create)
	mux.HandleFunc("DELETE /api/subscriptions/{id}", subH.Delete)

	notifyH := NewNotifyHandler(inventoryRepo, subRepo)
	mux.HandleFunc("POST /api/notify/inventory-alert", notifyH.InventoryAlert)

	return CORSMiddleware(mux)
}
