package main

import (
	"log"
	"net/http"
	"os"

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

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
