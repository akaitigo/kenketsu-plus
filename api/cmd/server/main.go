package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/handler"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	centerRepo := repository.NewCenterRepository()
	donationRepo := repository.NewDonationRepository()
	inventoryRepo := repository.NewInventoryRepository()
	subRepo := repository.NewSubscriptionRepository()
	calculator := service.NewDonationCalculator()

	router := handler.NewRouter(centerRepo, donationRepo, inventoryRepo, subRepo, calculator)

	// H-1: full timeout configuration
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Starting server on :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
